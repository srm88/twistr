package twistr

/*
[player] [verb] ...

[phase] [card] [play]

play:
op [realign [[country] roll/roll] ...] [opponent event ...]
op [coup [country] roll]
op [space roll]
op [influence N [country] ...]
event [event]
*discard [card] [roll]

********

olympics [participate [roll/roll]] | [boycott op ...]
summit [roll/roll] [defcon change]
howilearnedtostopworrying [defcon level]
junta [country] [coup | realign]

saltnegotiations [card]
fiveyearplan [card]
terrorism [card]x1 | 2
missileenvy [event] | [op [card] ... ]
grainsalestosoviets [card] [play ... ] | [return op ...]
asknotwhatyourcountry [discard replacement] ...
starwars [card]
latinamericandebtcrisis [card] | [country]x2
blockade [card] | 3US WGermany
aldrichamesremix [card]
ourmanintehran [card]x<5

warsawpactformed [remove [country]x4] | [add [country]x5]
socialistgovernments [country]x3
comecon [country]x4
trumandoctrine [country]
independentreds [country]
marshallplan [country]x7
suezcrisis [country]x4
easteuropeanunrest [country]x3
decolonization [country]x4
destalinization [country]x4 [country]x4
colonialrearguards [country]x4
puppetgovernments [country]x3
oasfounded [country]x2
thevoiceofamerica [country]x4
thereformer [country]x4 | 6
marinebarracksbombing [country]x2
pershingiideployed [country]x3
thecambridgefive [revealed scoring card]... [country]
specialrelationship [country]
southafricanunrest [south africa]x2 | [south africa] [adj to south africa]x2
muslimrevolution [country]x2
liberationtheology [country]x3
ussuririverskirmish nil | [country]x4

koreanwar [roll]
arabisraeliwar [roll]
indopakistaniwar [country] [roll]
brushwar [country] [roll]
iraniraqwar [country] [roll]

ciacreated [op ...]
unintervention [card] [op ...]
abmtreaty [op ... ]
lonegunman [op ... ]

sovietsshootdownkal007 [influence | realign]
glasnost [influence | realign]
teardownthiswall [coup | realign]
che [coup ...] && [coup ...]
ortegaelectedinnicaragua [coup ...]

chernobyl [region]

NORAD SPECIAL, EVERY ACTION ROUND [country]
*/
