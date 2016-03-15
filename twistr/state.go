package twistr

type State struct {
    VP int8

    Defcon uint8

    MilOps [2]uint8

    SpaceRace [2]uint8

    Turn uint8
    AR uint8

    Countries map[CountryId]*Country

    Events map[CardId]*Card

    Removed []*Card

    Discard []*Card

    Hands [2]map[CardId]*Card

    ChinaCardPlayer Aff
    ChinaCardFaceUp bool
}
