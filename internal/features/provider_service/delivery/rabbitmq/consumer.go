package rabbitmq

const (
	PublishedQueueName  = "provider_service_published"
	PublishedRoutingKey = "provider_service_published"

	UnpublishedQueueName  = "provider_service_unpublished"
	UnpublishedRoutingKey = "provider_service_unpublished"

	UpdatedQueueName  = "provider_service_updated"
	UpdatedRoutingKey = "provider_service_updated"

	MediaUploadedQueueName  = "media_uploaded"
	MediaUploadedRoutingKey = "media_uploaded"

	DeletedQueueName  = "provider_service_deleted"
	DeletedRoutingKey = "provider_service_deleted"
)
