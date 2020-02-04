package qdsservices

import (
	dtaf "github.com/theovassiliou/doctrans/dtaservice"
)

type QdsWormhole struct {
	dtaf.UnimplementedDTAServerServer

	// -- Galaxies Registrar -- The service has be registered there
	Register      bool   `opts:"group=Registrar" help:"Register service with EUREKA, if set"`
	RegistrarURL  string `opts:"group=Registrar" help:"Registry URL"`
	RegistrarUser string `opts:"group=Registrar" help:"Registry User, no user used if not provided"`
	RegistrarPWD  string `opts:"group=Registrar" help:"Registry User Password, no password used if not provided"`
	TTL           uint   `opts:"group=Registrar" help:"Time in seconds to reregister at Registrar."`

	// -- Resolver: A wormhole requires a resolver, because the wormwhole is registered in the outer galaxies but needs the resolver to
	// 		find the services in the inner galaxie. Remember: wormhole belong to two galaxies, innner and outer
	ResolverURL          string `opts:"group=Resolver" help:"Resolver URL"`
	ResolverUser         string `opts:"group=Resolver" help:"Resolver User, no user used if not provided"`
	ResolverPWD          string `opts:"group=Resolver" help:"Resolver User Password, no password used if not provided"`
	ResolverTTL          uint   `opts:"group=Resolver" help:"Time in seconds to reregister at Resolver."`
	ResolverRegistration bool   `opts:"group=Resolver" help:"Register in addition also to the resolver"`
}
