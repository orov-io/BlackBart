package server

import (
	"context"
	"io/ioutil"
	"os"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

const defaultFirebaseConfigName = "firebase.json"

// GetAuthClient returns the auth client attached to current service
func GetAuthClient() (*auth.Client, error) {
	service, err := GetService()
	if err != nil {
		GetLogger().
			WithError(err).
			Warn("Can't retrieve auth client. Does you call server.Init()??")
		return nil, err
	}

	authClient, err := service.GetAuthClient()
	if err != nil {
		GetLogger().
			WithError(err).
			Warn("Can't retrieve auth client. Does you add a firebase options??")
	}

	return authClient, err
}

// GetFirebaseApp returns the firebase app client
func GetFirebaseApp() (*firebase.App, error) {
	service, err := GetService()
	if err != nil {
		GetLogger().
			WithError(err).
			Warn("Can't retrieve the service. Does you call server.Init()??")
		return nil, err
	}

	firebaseApp, err := service.GetFirebaseApp()
	if err != nil {
		GetLogger().
			WithError(err).
			Warn("Can't retrieve firebase app. Does you add a firebase options??")
	}

	return firebaseApp, err
}

func (s *Service) initAuth() error {
	if !mustInitializeFirebase(s.options) {
		return NoFirebaseOptionsError()
	}

	credentials, err := s.getFirebaseCredentialsFromOptions()
	if err != nil {
		return err
	}
	s.firebase = initAuthApp(credentials)

	return nil
}

func mustInitializeFirebase(options *Options) bool {
	return options.firebase != nil
}

func (s *Service) getFirebaseCredentialsFromOptions() (credentials []byte, err error) {
	firebaseOptions := s.options.firebase
	if firebaseOptions.bucket != "" {
		return s.getCredentialsFromBucket(), nil
	}

	if firebaseOptions.configPath != "" {
		return s.getCredentialsFromFile(firebaseOptions.configPath), nil
	}
	return nil, NoFirebaseOptionsError()
}

func (s *Service) getCredentialsFromBucket() []byte {
	ctx := context.Background()
	bucket := s.options.firebase.bucket
	name := getFirebaseConfigFileNameFromOptions(s.options.firebase)

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		GetLogger().WithError(err).Fatalf("Can't get storage connection")
	}
	credentials, err := storageClient.Bucket(bucket).Object(name).NewReader(ctx)
	if err != nil {
		GetLogger().WithError(err).Fatal("Can't stablish reader connection")
	}

	buffer, err := ioutil.ReadAll(credentials)
	if err != nil {
		logrus.Fatalf("Can't read credentials: %v", err)
	}

	return buffer
}

func getFirebaseConfigFileNameFromOptions(options *FirebaseOptions) string {
	name := options.name
	if options.name == "" {
		name = defaultFirebaseConfigName
	}

	return name
}

func (s *Service) getCredentialsFromFile(filePath string) []byte {
	file, err := os.Open(filePath)
	if err != nil {
		GetLogger().WithError(err).Fatalf("Can't open credentials file")
	}

	fileStats, err := file.Stat()
	if err != nil {
		GetLogger().WithError(err).Fatalf("Can't stat credentials file")
	}

	buffer := make([]byte, fileStats.Size())
	_, err = file.Read(buffer)

	if err != nil {
		GetLogger().WithError(err).Fatalf("Can't read credentials file")
	}
	return buffer
}

// initAuthApp starts a firebase app with provided credentials
func initAuthApp(firebaseCredentials []byte) *firebase.App {
	context := context.Background()
	opt := option.WithCredentialsJSON(firebaseCredentials)
	firebaseApp, err := firebase.NewApp(context, nil, opt)
	if err != nil {
		GetLogger().WithError(err).Fatal("error initializing app:")
	}
	return firebaseApp
}
