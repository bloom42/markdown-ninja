package service

import (
	"fmt"
	"text/template"

	"github.com/skerkour/stdx-go/db"
	"github.com/skerkour/stdx-go/queue"
	"markdown.ninja/cmd/mdninja-server/config"
	"markdown.ninja/pingoo-go"
	"markdown.ninja/pkg/mailer"
	"markdown.ninja/pkg/services/contacts"
	"markdown.ninja/pkg/services/content"
	"markdown.ninja/pkg/services/emails"
	"markdown.ninja/pkg/services/events"
	"markdown.ninja/pkg/services/kernel"
	"markdown.ninja/pkg/services/organizations"
	"markdown.ninja/pkg/services/store/notifications"
	"markdown.ninja/pkg/services/store/repository"
	"markdown.ninja/pkg/services/websites"
)

type StoreService struct {
	repo   repository.StoreRepository
	db     db.DB
	queue  queue.Queue
	mailer mailer.Mailer

	websitesService      websites.Service
	kernel               kernel.PrivateService
	contentService       content.Service
	contactsService      contacts.Service
	eventsService        events.Service
	emailsService        emails.Service
	organizationsService organizations.Service
	pingoo               *pingoo.Client

	httpConfig                     config.Http
	websitesPort                   string
	orderConfirmationEmailTemplate *template.Template
}

func NewStoreService(db db.DB, queue queue.Queue, conf config.Config, mailer mailer.Mailer, kernel kernel.PrivateService, websitesService websites.Service,
	contentService content.Service, contactsService contacts.Service, eventsService events.Service,
	emailsService emails.Service, organizationsService organizations.Service, pingoo *pingoo.Client) (service *StoreService, err error) {
	repo := repository.NewStoreRepository()

	orderConfirmationEmailTemplate, err := template.New("store.OrderConfirmationEmailTemplate").Parse(notifications.OrderConfirmationEmailTemplate)
	if err != nil {
		err = fmt.Errorf("store.NewService: Parsing orderConfirmationEmailTemplate: %w", err)
		return
	}

	service = &StoreService{
		repo:   repo,
		db:     db,
		queue:  queue,
		mailer: mailer,

		kernel:               kernel,
		websitesService:      websitesService,
		contentService:       contentService,
		contactsService:      contactsService,
		eventsService:        eventsService,
		emailsService:        emailsService,
		organizationsService: organizationsService,
		pingoo:               pingoo,

		httpConfig:                     conf.HTTP,
		websitesPort:                   conf.HTTP.WebsitesPort,
		orderConfirmationEmailTemplate: orderConfirmationEmailTemplate,
	}
	return
}
