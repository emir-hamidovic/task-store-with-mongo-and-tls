module rest

go 1.16

require (
	go.mongodb.org/mongo-driver v1.8.4 // indirect
	taskstore v1.0.0
)

replace taskstore v1.0.0 => ./taskstore
