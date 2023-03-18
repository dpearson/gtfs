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

	// RouteTypeTrolleybus indicates that the route is a trolleybus route.
	RouteTypeTrolleybus

	// RouteTypeMonorail indicates that the route is a monorail route.
	RouteTypeMonorail

	// EXTENDED ROUTE TYPES
	//
	// The following types are extended route types proposed by Google, although
	// they are not part of the GTFS standard.
	//
	// See: https://developers.google.com/transit/gtfs/reference/extended-route-types

	RouteTypeExtendedRailwayService
	RouteTypeExtendedHighSpeedRail
	RouteTypeExtendedLongDistanceRail
	RouteTypeExtendedInterRegionalRail
	RouteTypeExtendedCarTransportRail
	RouteTypeExtendedSleeperRail
	RouteTypeExtendedRegionalRail
	RouteTypeExtendedTouristRail
	RouteTypeExtendedRailShuttle
	RouteTypeExtendedSuburbanRail
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
	RouteTypeExtendedMetro
	RouteTypeExtendedUnderground
	RouteTypeExtendedUrbanRail
	RouteTypeExtendedAllUrbanRail
	RouteTypeExtendedMonorail

	RouteTypeExtendedBusService
	RouteTypeExtendedRegionalBusService
	RouteTypeExtendedExpressBusService
	RouteTypeExtendedStoppingBusService
	RouteTypeExtendedLocalBusService
	RouteTypeExtendedNightBusService
	RouteTypeExtendedPostBusService
	RouteTypeExtendedSpecialNeedsBusService
	RouteTypeExtendedMobilityBusService
	RouteTypeExtendedMobilityBusForRegisteredDisabledService
	RouteTypeExtendedSightseeingBus
	RouteTypeExtendedShuttleBus
	RouteTypeExtendedSchoolBus
	RouteTypeExtendedSchoolAndPublicServiceBus
	RouteTypeExtendedRailReplacementBusService
	RouteTypeExtendedDemandAndResponseBusService
	RouteTypeExtendedAllBusServices

	RouteTypeExtendedTrolleybusService

	RouteTypeExtendedTramService
	RouteTypeExtendedCityTramService
	RouteTypeExtendedLocalTramService
	RouteTypeExtendedRegionalTramService
	RouteTypeExtendedSightseeingTramService
	RouteTypeExtendedShuttleTramService
	RouteTypeExtendedAllTramServices

	RouteTypeExtendedWaterTransportService

	RouteTypeExtendedAirService

	RouteTypeExtendedFerryService

	RouteTypeExtendedAerialLiftService
	RouteTypeExtendedTelecabinService
	RouteTypeExtendedCableCarService
	RouteTypeExtendedElevatorService
	RouteTypeExtendedChairLiftService
	RouteTypeExtendedDragLiftService
	RouteTypeExtendedSmallTelecabinService
	RouteTypeExtendedAllTelecabinServices

	RouteTypeExtendedFunicularService

	RouteTypeExtendedTaxiService
	RouteTypeExtendedCommunalTaxiService
	RouteTypeExtendedWaterTaxiService
	RouteTypeExtendedRailTaxiService
	RouteTypeExtendedBikeTaxiService
	RouteTypeExtendedLicensedTaxiService
	RouteTypeExtendedPrivateHireServiceVehicle
	RouteTypeExtendedAllTaxiServices

	RouteTypeExtendedMiscellaneousService
	RouteTypeExtendedHorseDrawnCarriage
)

var routeTypes = map[string]RouteType{
	"0":  RouteTypeLightRail,
	"1":  RouteTypeSubway,
	"2":  RouteTypeRail,
	"3":  RouteTypeBus,
	"4":  RouteTypeFerry,
	"5":  RouteTypeCableCar,
	"6":  RouteTypeGondola,
	"7":  RouteTypeFunicular,
	"11": RouteTypeTrolleybus,
	"12": RouteTypeMonorail,

	// See: https://developers.google.com/transit/gtfs/reference/extended-route-types
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

	"400": RouteTypeExtendedUrbanRailService,
	"401": RouteTypeExtendedMetro,
	"402": RouteTypeExtendedUnderground,
	"403": RouteTypeExtendedUrbanRail,
	"404": RouteTypeExtendedAllUrbanRail,
	"405": RouteTypeExtendedMonorail,

	"700": RouteTypeExtendedBusService,
	"701": RouteTypeExtendedRegionalBusService,
	"702": RouteTypeExtendedExpressBusService,
	"703": RouteTypeExtendedStoppingBusService,
	"704": RouteTypeExtendedLocalBusService,
	"705": RouteTypeExtendedNightBusService,
	"706": RouteTypeExtendedPostBusService,
	"707": RouteTypeExtendedSpecialNeedsBusService,
	"708": RouteTypeExtendedMobilityBusService,
	"709": RouteTypeExtendedMobilityBusForRegisteredDisabledService,
	"710": RouteTypeExtendedSightseeingBus,
	"711": RouteTypeExtendedShuttleBus,
	"712": RouteTypeExtendedSchoolBus,
	"713": RouteTypeExtendedSchoolAndPublicServiceBus,
	"714": RouteTypeExtendedRailReplacementBusService,
	"715": RouteTypeExtendedDemandAndResponseBusService,
	"716": RouteTypeExtendedAllBusServices,

	"800": RouteTypeExtendedTrolleybusService,

	"900": RouteTypeExtendedTramService,
	"901": RouteTypeExtendedCityTramService,
	"902": RouteTypeExtendedLocalTramService,
	"903": RouteTypeExtendedRegionalTramService,
	"904": RouteTypeExtendedSightseeingTramService,
	"905": RouteTypeExtendedShuttleTramService,
	"906": RouteTypeExtendedAllTramServices,

	"1000": RouteTypeExtendedWaterTransportService,

	"1100": RouteTypeExtendedAirService,

	"1200": RouteTypeExtendedFerryService,

	"1300": RouteTypeExtendedAerialLiftService,
	"1301": RouteTypeExtendedTelecabinService,
	"1302": RouteTypeExtendedCableCarService,
	"1303": RouteTypeExtendedElevatorService,
	"1304": RouteTypeExtendedChairLiftService,
	"1305": RouteTypeExtendedDragLiftService,
	"1306": RouteTypeExtendedSmallTelecabinService,
	"1307": RouteTypeExtendedAllTelecabinServices,

	"1400": RouteTypeExtendedFunicularService,

	"1500": RouteTypeExtendedTaxiService,
	"1501": RouteTypeExtendedCommunalTaxiService,
	"1502": RouteTypeExtendedWaterTaxiService,
	"1503": RouteTypeExtendedRailTaxiService,
	"1504": RouteTypeExtendedBikeTaxiService,
	"1505": RouteTypeExtendedLicensedTaxiService,
	"1506": RouteTypeExtendedPrivateHireServiceVehicle,
	"1507": RouteTypeExtendedAllTaxiServices,

	"1700": RouteTypeExtendedMiscellaneousService,
	"1702": RouteTypeExtendedHorseDrawnCarriage,
}

func parseRouteType(val string) (RouteType, error) {
	routeType, ok := routeTypes[val]
	if !ok {
		return 0, fmt.Errorf("invalid route type: %s", val)
	}

	return routeType, nil
}
