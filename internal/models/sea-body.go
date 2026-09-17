package models

type Secretariat string

const (
	SecMedia       Secretariat = "media"
	SecAcademic    Secretariat = "academic"
	SecSports      Secretariat = "sports"
	SecExRelations Secretariat = "external_relations"
	SecCultural    Secretariat = "cultural"
	SecFinancial   Secretariat = "financial"
	SecGeneral     Secretariat = "general"
	SecSocial      Secretariat = "social"
)

type Department string

const (
	DEP_MECHANICAL   Department = "mechanical"
	DEP_CIVIL        Department = "civil"
	DEP_ELECTRICAL   Department = "electrical"
	DEP_CHEMICAL     Department = "chemical"
	DEP_PETROLEUM    Department = "petroleum"
	DEP_AGRICULTURAL Department = "agricultural"
	DEP_MINING       Department = "mining"
	DEP_SURVEYING    Department = "surveying"
)
