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

olympics [participate [roll/roll]] | [op ...]
summit [roll/roll] [defcon change]
how-i-learned-to-stop-worrying [defcon level]
junta [country] [coup | realign]

salt-negotiations [card]
five-year-plan [card]
terrorism [card]x1 | 2
missile-envy [event] | [op [card] ... ]
grain-sales-to-soviets [card] [play ... ] | [return op ...]
ask-not-what-your-country [discard replacement] ...
star-wars [card]
latin-american-debt-crisis [card] | [country]x2
blockade [card] | -3US WGermany
aldrich-ames-remix [card]
our-man-in-tehran [card]x<5

warsaw-pact-formed [remove [country]x4] | [add [country]x5]
socialist-governments [country]x3
comecon [country]x4
truman-doctrine [country]
independent-reds [country]
marshall-plan [country]x7
suez-crisis [country]x4
east-european-unrest [country]x3
decolonization [country]x4
de-stalinization [country]x4 [country]x4
colonial-rear-guards [country]x4
puppet-governments [country]x3
oas-founded [country]x2
the-voice-of-america [country]x4
the-reformer [country]x4 | 6
marine-barracks-bombing [country]x2
pershing-ii-deployed [country]x3
the-cambridge-five [revealed scoring card]... [country]
special-relationship [country]

korean-war [roll]
arab-israeli-war [roll]
indo-pakistani-war [country] [roll]
brush-war [country] [roll]
iran-iraq-war [country] [roll]

cia-created [op ...]
un-intervention [card] [op ...]
abm-treaty [op ... ]

soviets-shoot-down-kal-007 [influence | realign]
glasnost [influence | realign]
tear-down-this-wall [coup | realign]
che [coup ...] && [coup ...]

chernobyl [region]

NORAD SPECIAL, EVERY ACTION ROUND [country]
*/
