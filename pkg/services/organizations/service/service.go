package service

import (
	"html/template"

	"github.com/bloom42/stdx-go/db"
	"github.com/bloom42/stdx-go/queue"
	"markdown.ninja/cmd/mdninja-server/config"
	"markdown.ninja/pingoo-go"
	"markdown.ninja/pkg/mailer"
	"markdown.ninja/pkg/services/content"
	"markdown.ninja/pkg/services/events"
	"markdown.ninja/pkg/services/kernel"
	"markdown.ninja/pkg/services/organizations/repository"
	"markdown.ninja/pkg/services/organizations/templates"
	"markdown.ninja/pkg/services/store"
	"markdown.ninja/pkg/services/websites"
)

type OrganizationsService struct {
	repo         repository.OrganizationsRepository
	db           db.DB
	mailer       mailer.Mailer
	queue        queue.Queue
	stripeConfig config.Stripe
	httpConfig   config.Http
	isSelfHosted bool

	kernel          kernel.PrivateService
	websitesService websites.Service
	eventsService   events.Service
	contentService  content.Service
	storeService    store.Service
	pingoo          *pingoo.Client

	staffInvitationEmailTemplate *template.Template
}

func NewOrganizationsService(conf config.Config, db db.DB, mailer mailer.Mailer, queue queue.Queue,
	kernel kernel.PrivateService, pingoo *pingoo.Client) *OrganizationsService {
	repo := repository.NewOrganizationsRepository()

	staffInvitationEmailTemplate := template.Must(template.New("organizations.StaffInvitationEmailTemplate").Parse(templates.StaffInvitationEmailTemplate))

	return &OrganizationsService{
		repo:         repo,
		db:           db,
		mailer:       mailer,
		queue:        queue,
		stripeConfig: *conf.Stripe,
		httpConfig:   conf.HTTP,
		isSelfHosted: !conf.Saas,

		kernel:          kernel,
		websitesService: nil,
		eventsService:   nil,
		contentService:  nil,
		storeService:    nil,

		pingoo: pingoo,

		staffInvitationEmailTemplate: staffInvitationEmailTemplate,
	}
}

func (service *OrganizationsService) InjectServices(websitesService websites.Service, eventsService events.Service,
	contentService content.Service, storeService store.Service) {
	service.websitesService = websitesService
	service.eventsService = eventsService
	service.contentService = contentService
	service.storeService = storeService
}
