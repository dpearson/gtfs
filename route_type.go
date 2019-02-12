package gtfs

import "fmt"

// RouteType specifies the type of vehicles operating on a route.
type RouteType int

const (
	// RouteTypeNotSpecified indicates that no route type was specified.
	RouteTypeNotSpecified RouteType = iota

	// RouteTypeLightRail indicates that the route is a light rail route.
	RouteTypeLightRail

	// RouteTypeSubway indicates that the route is a subway or metro route.
	RouteTypeSubway

	// RouteTypeRail indicates that the route is an intercity/long-distance rail
	// route.
	RouteTypeRail

	// RouteTypeBus indicates that the route is a bus route.
	RouteTypeBus

	// RouteTypeFerry indicates that the route is a ferry route.
	RouteTypeFerry

	// RouteTypeCableCar indicates that the route is a cable car route.
	RouteTypeCableCar

	// RouteTypeGondola indicates that the route is an aerial gondola route.
	RouteTypeGondola

	// RouteTypeFunicular indicates that the route is a funicular route.
	RouteTypeFunicular

	// EXTENDED ROUTE TYPES
	//
	// The following types are extended route types proposed by Google, although
	// they are not part of the GTFS standard.

	RouteTypeExtendedRailwayService
	RouteTypeExtendedHighSpeedRail
	RouteTypeExtendedLongDistanceRail
	RouteTypeExtendedInterRegionalRail
	RouteTypeExtendedCarTransportRail
	RouteTypeExtendedSleeperRail
	RouteTypeExtendedRegionalRail
	RouteTypeExtendedTouristRail
	RouteTypeExtendedRailShuttle
	RouteTypeExtendedSuburbanRail // used for route type codes 109 and 300
	RouteTypeExtendedReplacementRail
	RouteTypeExtendedSpecialRail
	RouteTypeExtendedLorryTransportRail
	RouteTypeExtendedAllRail
	RouteTypeExtendedCrossCountryRail
	RouteTypeExtendedVehicleTransportRail
	RouteTypeExtendedRackAndPinionRail
	RouteTypeExtendedAdditionalRail

	RouteTypeExtendedCoachService
	RouteTypeExtendedInternationalCoach
	RouteTypeExtendedNationalCoach
	RouteTypeExtendedShuttleCoach
	RouteTypeExtendedRegionalCoach
	RouteTypeExtendedSpecialCoach
	RouteTypeExtendedSightseeingCoach
	RouteTypeExtendedTouristCoach
	RouteTypeExtendedCommuterCoach
	RouteTypeExtendedAllCoach

	RouteTypeExtendedUrbanRailService
	RouteTypeExtendedMetro       // used for route type codes 401 and 500
	RouteTypeExtendedUnderground // used for route type codes 402 and 600
	RouteTypeExtendedUrbanRail
	RouteTypeExtendedAllUrbanRail
	RouteTypeExtendedMonorail
)

var routeTypes = map[string]RouteType{
	"0":   RouteTypeLightRail,
	"1":   RouteTypeSubway,
	"2":   RouteTypeRail,
	"3":   RouteTypeBus,
	"4":   RouteTypeFerry,
	"5":   RouteTypeCableCar,
	"6":   RouteTypeGondola,
	"7":   RouteTypeFunicular,
	"100": RouteTypeExtendedRailwayService,
	"101": RouteTypeExtendedHighSpeedRail,
	"102": RouteTypeExtendedLongDistanceRail,
	"103": RouteTypeExtendedInterRegionalRail,
	"104": RouteTypeExtendedCarTransportRail,
	"105": RouteTypeExtendedSleeperRail,
	"106": RouteTypeExtendedRegionalRail,
	"107": RouteTypeExtendedTouristRail,
	"108": RouteTypeExtendedRailShuttle,
	"109": RouteTypeExtendedSuburbanRail,
	"110": RouteTypeExtendedReplacementRail,
	"111": RouteTypeExtendedSpecialRail,
	"112": RouteTypeExtendedLorryTransportRail,
	"113": RouteTypeExtendedAllRail,
	"114": RouteTypeExtendedCrossCountryRail,
	"115": RouteTypeExtendedVehicleTransportRail,
	"116": RouteTypeExtendedRackAndPinionRail,
	"117": RouteTypeExtendedAdditionalRail,
	"200": RouteTypeExtendedCoachService,
	"201": RouteTypeExtendedInternationalCoach,
	"202": RouteTypeExtendedNationalCoach,
	"203": RouteTypeExtendedShuttleCoach,
	"204": RouteTypeExtendedRegionalCoach,
	"205": RouteTypeExtendedSpecialCoach,
	"206": RouteTypeExtendedSightseeingCoach,
	"207": RouteTypeExtendedTouristCoach,
	"208": RouteTypeExtendedCommuterCoach,
	"209": RouteTypeExtendedAllCoach,
	"300": RouteTypeExtendedSuburbanRail,
	"400": RouteTypeExtendedUrbanRailService,
	"401": RouteTypeExtendedMetro,
	"402": RouteTypeExtendedUnderground,
	"403": RouteTypeExtendedUrbanRail,
	"404": RouteTypeExtendedAllUrbanRail,
	"405": RouteTypeExtendedMonorail,
	"500": RouteTypeExtendedMetro,
	"600": RouteTypeExtendedUnderground,
}

func parseRouteType(val string) (RouteType, error) {
	routeType, ok := routeTypes[val]
	if !ok {
		return 0, fmt.Errorf("invalid route type: %s", val)
	}

	return routeType, nil
}
