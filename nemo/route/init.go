package route

import "fmt"

var ENTENDU_EN_VOYAGE = "/entendu_en_voyage"
var RENCONTRE_EN_VOYAGE = "/rencontre_en_voyage"
var CLEF_CANONIQUE = "/clefs_canoniques"
var VAISEAU = "/vaiseau"
var REMERCIEMENTS = "/remerciements"
var ACCOSTE_EN_VOYAGE_ROOT = "/accoste_en_voyage"
var ACCOSTE_EN_VOYAGE = fmt.Sprintf("%v/{corpus}/{idproprio}", ACCOSTE_EN_VOYAGE_ROOT)
