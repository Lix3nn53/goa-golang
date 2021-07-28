# Guardians of Adelia REST API

I am studiying golang and rest api best practices while developing backend for [Guardians of Adelia React App](https://github.com/Lix3nn53/goa-react-app).
This web application is for [Guardians of Adelia](https://github.com/Lix3nn53/GuardiansOfAdelia) which is a dungeon driven mmorpg minecraft server project.

## Web Framework

I selected [gin](https://github.com/gin-gonic/gin) for framework because of simplicity and popularity since I am new to golang.

## Docker

My first time using docker. Use docker-compose with api and proxy service. Optimize image size of api with multi-stage build.

## Documentation

Use [go-swagger](https://github.com/go-swagger/go-swagger) to create documentation and serve with [Redocly/redoc](https://github.com/Redocly/redoc)

## Repository-Service Pattern with Dependency Injection

I studied dependency injection(DI) and implemented it with repository-service pattern. Dependency Injection (DI) is a design pattern used to implement IoC. It allows the creation of dependent objects outside of a class and provides those objects to a class through different ways. Using DI, we move the creation and binding of the dependent objects outside of the class that depends on them.

The Dependency Injection pattern involves 3 types of classes.

- Client Class: The client class (dependent class) is a class which depends on the service class.
- Service Class: The service class (dependency) is a class that provides service to the client class.
- Injector Class: The injector class injects the service class object into the client class.

I used [google/wire](https://github.com/google/wire) to create container because of simplicity, it generates obvious and readable code.

Got a lot of help from [jaumeCloquellCapo/go-api-boirplate](https://github.com/jaumeCloquellCapo/go-api-boirplate) while implementing repository-service pattern.

## Testing

When testing a component, we ideally want to isolate it completely but this is especially harder when the component we want to test has dependencies on other components from different layers in our software. To promote the desired isolation, it is common for developers to write fake simplified implementations of those dependencies to be used during the tests. Generate rates mock interfaces from a source file with [gomock](https://github.com/golang/mock).

## Extra Notes

### Logger

Used [uber-go/zap](https://github.com/uber-go/zap) for custom logger.
