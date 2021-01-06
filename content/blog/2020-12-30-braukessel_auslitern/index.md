+++
draft = false
toc = true
author = "Till Klocke"
removeBlur = false
date = "2021-01-06T14:20:00+01:00"
comments = false
title = "Dein Braukessel belügt dich vielleicht"
description = "Braukessel auslitern und eigene Skala ätzen"
tags = ["Brauheld", "Malzrohrsysteme", "Ausrüstung", "Klarstein", "HowTo"]
categories = ["Anleitung"]
publishdate = "2021-01-06T14:20:00+01:00"
[[images]]
  src = "/blog/2020-12-30-braukessel_auslitern/IMG_20201219_141627.jpg"
  alt = "geätzte Literskala"
  stretch = "c"
+++
{{< toc >}}  

Wie bereits erwähnt verwende ich den Klartstein Brauheld 30 um meine Würze herzustellen. Und wie
ebenfalls bereits erwähnt ist das nicht unbedingt ein hochwertiges Gerät (auch wenn es grundsätzlich
funktioniert). Das betrifft auch die eingeprägte Literskala. Diese ist verschoben und verzerrt wodurch
keinerlei sinnvollen volumetrischen Ablesungen möglich sind. Weder absolute Werte, noch Differenzwerte
(wie z.B. die Verdampfung) können von dieser Skala abgelesen werden.

{{< img "IMG_20201219_141627.jpg" "Geprägte Literskala neben der neu geätzten" >}}

## Überführe deinen Kessel der Lüge

Natürlich kann auch mit einer solch falschen Literskala gebraut werden, aber es wird schnell offensichtlich,
dass berechnete Rezepte dann wenig Sinn machen bzw. der Nutzen vom exakten Berechnen von Rezepten gering ist.
Es ist also nützlich ggf. bereits vor dem ersten Brauen zu überprüfen, ob die Literskala im Braukessel korrekt ist.

Hierzu sollte ein ausreichend großes Gefäß und eine Waage vorhanden sein. Auch wenn das Gefäß eine
volumetrische Skala hat, würde ich dieser nicht unbedingt mehr trauen als der Skala im Braukessel. Meiner
Erfahrung nach sind Küchenwaagen präziser als volumetrischen Skalen auf Küchengeräten. In meinem Fall
habe ich eine Küchenwaage und einen 3L Messbecher verwendet. Da Wasser ja recht exakt 1000g bei 1L Volumen wiegt,
habe ich immer ausreichend Wasser bis zu einer Markierung abgewogen. Liegt die erste Markierung also bei 5L,
werden insgesamt 5000g Wasser abgewogen und in den Kessel gegeben. Das kann natürlich auch in kleinere Gaben, 
z.B. á 1000g geschehen.
Ist die erste Markierung bei 5L und es sind 5000g Wasser im Kessel, sollte der Wasserspiegel
nun bei dieser Markierung sein. Ist dies nicht der Fall, kann hier bereits abgebrochen werden, da es recht sicher
ist dass die Literskala in deinem Braukessel lügt.

Aber selbst, wenn die erste Markierung korrekt aussieht, lohnt es sich genügend Wasser abzuwiegen um auch
die nächste Markierung zu prüfen. Also z.B. weitere 5000g Wasser abzuwiegen um dann 10.000g Wasser im
Braukessel zu haben, was dann den Wasserspiegel bis zur 10L Markierung heben sollte. Wenn auch dies passt,
ist es recht sicher, dass die Skala im Braukessel korrekt ist. Sollte die persönliche Paranoidität eher größer sein, 
kann dieser Prozess natürlich auch bis zur letzten Markierung wiederholt werden.

## Mein Kessel belügt mich, was kann ich tun?

Sollte die Literskala nicht korrekt sein, muss eine eigene Literskala her. Das ist tatsächlich gar nicht so schwierig, 
sogar ich habe das halbwegs zufriedenstellend hinbekommen. Alles was dazu benötigt wird:

* Eine Küchenwaage
* Ein ausreichend großes Gefäß (zwischen 1L und 5L)
* Klebeband. Es sollte auch auf feuchtem Stahl gut kleben und einen guten Kontrast zu Edelstahl bilden (also nicht unbedingt graues)
* Schablonen für Striche und Zahlen
* 9V Batterie
* Wattestäbchen
* Kabel
* Glasschüssel
* 5-6% Weißweinessig
* Salz

### But first: Safety!

Dieser ganze Prozess ist zwar nicht sonderlich gefährlich, und es wird auch nicht mit allzu
gefährlichen Stoffen hantiert. Trotzdem sollten ein paar Kleinigkeiten beachtet werden. Beim Elektroätzen
entstehen Gase. Führt diesen ganzen Prozess also nur an einem gut belüfteten Ort oder am besten
draußen durch.
Auch wenn es sich nur um Weißweinessig und Salt handelt, habe ich Nitrilhandschuhe und eine
Schutzbrille getragen. Das empfehle ich alleine schon um eine Gewisse Gewohnheit bei der Einhaltung
von Schutzmaßnahmen zu fördern. Weißweinessig gemischt mit Salz ist zumindest sehr unangenehm im Auge.

Außerdem solltet ihr auch euren Braukessel schützen. Damit meine ich, dass ihr diesen Prozess erst einmal
an etwas anderem testen solltet. Versucht einen alten Growler, Tasse oder ähnliches aus Edelstahl zu
markieren und bekommt ein Gefühl dafür nach welcher Zeit ihr welche Markierung erreicht und wie schnell
oder langsam ihr das Wattestäbchen führen solltet. Da wir hier nicht kalibrierten Werkzeugen arbeiten,
ist es kaum möglich vorzugeben wie lange ihr für eine ausreichende Markierung braucht und wie schnell ihr
das Wattestäbchen führen solltet.

### Schablonen

Für die Markierung müssen passende Schablonen gesucht oder selbst gebastelt werden. Es gibt in Kreativmärkten
etc. Schablonen mit Zahlen und Strichen. Ich habe aber nichts gefunden was mir gefallen hätte. Entweder
sagte mir die Schriftart nicht zu, oder die Zahlen waren schlicht zu groß. Daher habe ich in Fusion360
eigene {{< download "LiterScaleStencils.f3d" >}}Schablonen{{< /download >}} designed. Die Zahl auf der Schablone
kann in der Skizze einfach angepasst werden.
{{< img "IMG_20200818_185716.jpg" "3D-gedruckte Schablonen" >}}
So habe ich mehrere Schablonen generiert und zusätzlich eine, mit einem kürzeren Strich für die Liter.
Beim erstellen von Schablonen muss darauf geachtet werden, dass die Schriftart keine losen Teile erzeugt.
Also z.B. bei der "0" der innere Teil noch Brücken zum äußeren Teil hat. Es gibt Schriftarten, die so aufgebaut
sind. Bei mir hat sich die Schriftart "Octin Stencil" bewährt.
Die Schablonen habe ich mit meinem 3D-Drucker gedruckt. Dabei muss sichergestellt sein, dass die Schablonen flexibel
genug bleiben um fest and die Kesselwand gedrückt werden zu können. Ich habe mich für 0,5mm Dicke mit bei einer
Schichthöhe von 0,1mm in PLA entschieden. Das hat soweit ganz gut geklappt.
Alternativ kann so eine Schablone aber vermutlich auch mittels eines feinen Messers aus Kunststoffolie o.ä.
ausgeschnitten werden.

### Auslitern

{{< img "IMG_20200819_182618.jpg" "Klebeband zum Auslitern an der Kesselwand" >}}

Als erstes habe ich den Kessel gründlich sauber gemacht und in dem Bereich, in der die Skala sein soll,
auch mit Isopropanol-Alkohol gereinigt, damit das Klebeband auf jeden Fall gut hält.
Dann habe ich einen langen Streifen Klebeband vertikal an der Kesselwand von ganz unten nach ganz oben
angebracht. Dieser Streifen dient zur Orientierung damit die neue Skala bündig ausgerichtet ist.
Ihr könnt auf dem Bild sehen, dass ich meine Schablonen an diesem Streifen ausgerichtet habe und die
horizontalen Klebestreifen bei jeder Markierung dazu benutzt habe die Position meiner Schablonen zu markieren
und nicht den Wasserspiegel. Das hat den Vorteil, dass die Klebestreifen ein wenig über dem Wasserspiegel
angebracht werden können, was imho einfacher ist.
Die Schablonen habe ich immer temporär zwischen Zwei Magneten befestigt. Bei selbstklebenden Schablonen ist das
natürlich nicht nötig.
Ich habe mich damals dafür entschieden, dass ich nur alle 5L eine Zahl auf der Skala haben möchte und
dazwischen nur Striche. Daher habe ich auch nur in 5L-Inkrementen ausgelitert. Die kleinere 1L-Schritte
habe ich später einfach per Zentimetermaß ausgemessen. 

### Strom, Säure, Spaß

Da wir nun die Markierungen für die Position unserer Schablonen haben (oder die Schablonen direkt aufgeklebt
haben), können wir nun markieren. Allerdings ist Edelstahl ganz bewusst ein ziemlich widerstandsfähiges Material.
Immerhin muss es nicht nur Maische und Würze, sondern auch aggressiven Reinigungsmitteln, widerstehen. Daher müssen wir
ein wenig tricksen um Edelstahl so zu bearbeiten, dass eine visuelle Veränderung wahrgenommen werden kann.
Wir benutzen dafür mit Weißweinessig eine schwache Säure und verstärken quasi ihre Wirkung mit Strom. Die chemischen
Hintergründe, warum das funktioniert, warum Gleichspannung eher helle Markierungen und Wechelspannung eher dunkle
Markierungen hinterlässt, sind sehr spannend. Leider würde das den Rahmen hier aber sprengen. Wenn ihr aber einen
Chemie- oder Wissenschaftsnerd kennt, fragt und ich wette ihr werdet einen spannenden Vortrag bekommen.

Zunächst rühren wir uns eine Mischung aus dem Weißweinessig und dem Salz an. Das Salz wird nur dafür benötigt den
Essig besser leitend zu machen, damit Strom fließen kann. Daher spielt die Dosierung eine eher untergeordnete Rolle.
Ich habe einfach ca. 150ml Weißweinessig mit einem Esslöffel Salz gemischt. 150ml ist übrigens viel mehr als benötigt
wird. Es reicht auch bedeutend weniger. Lediglich das Wattestäbchen muss ausreichend feucht gehalten werden. 
{{< img "etching.png" "Schematische Darstellung des Ätzaufbaus" >}}
Nun kann der Plus-Pol der Batterie mit dem Kessel verbunden werden. Dafür habe ich einfach ein altes Stück Kabel und
etwas Klebeband verwendet. Der Minus-Pol der Batterie wird dem Wattestäbchen verbunden. Dabei ist es wichtig, dass
der unisolierte Teil des Kabels Kontakt zu dem saugfähigen Ende des Wattestäbchens hat. Dazu habe ich, 
beginnend bei dem saugfähigen Ende, einfach den entisolierten Teil des Kabels mit mehreren Windungen um das 
Wattestäbchen gewickelt und dann mit Klebeband fixiert. Der entisolierte Teil sollte aber auf keinen Fall den ganzen
Wattekopf umwickeln, damit das Kabel nie den Topf selbst berührt. Der Strom soll immer durch die Weinessig-Salz-Lösung
fließen. Bei dem Kabel am Wattestäbchen sollte natürlich darauf geachtet werden, dass es lang genug ist.

Wurde alles so vorbereitet und die Schablone an der gewünschten Stelle fixiert, wird das Wattestäbchen
mit der Weißwenigessig-Salz-Lösung befeuchtet. Dann wird das befeuchtete Wattestäbchen gleichmäßig über die
Ausschnitte in der Schablone gestrichen. Dieser Prozess dauert ein wenig. Es ist wichtig darauf zu achten, 
das Wattestäbchen kontinuierlich und gleichmäßig zu bewegen. Nur so wird die Markierung gleichmäßig und deutlich. 
Auch ist es wichtig, dass die Schablone wirklich bündig an der Kesselwand anliegt. Andernfalls besteht die Gefahr, 
dass die Markierung unscharf wird. Ein leises Zischen und Blubbern sind gute Zeichen, dass alles funktioniert. 
Nach und nach werden so alle Positionen mit Hilfe der Schablonen markiert.

Sollten, so wie bei mir, nur größere Abstände markiert worden sein, lohnt es sich nun die Abstände zwischen
den großen Markierungen zu messen. Ich habe das zwischen mehreren von den 5L-Markierungen gemacht. Die Abstände
sollten immer möglichst gleich sein. So kann auch geprüft werden wie exakt vorher gearbeitet wurde. Wurde
der Abstand zwischen den großen Markierungen ermittelt, wird dieser einfach durch die Zahl der kleinen
Markierungen geteilt. Sind also bisher alle 5L markiert worden und sollen Litermarkierungen hinzugefügt werden,
wird der Abstand einfach durch 5 geteilt. Nun können die kleineren Markierungen z.B. mit einem Stift oder
wieder mit Klebeband markiert werden, damit später die Schablone für die kleineren Markierungen angelegt werden kann.
Zuletzt werden die kleinen Markierungen genauso wie die großen geätzt.

## Ende und aus

{{< img "IMG_20201219_141627.jpg" "Selbstgeätzte Literskala neben der geprägten" >}}
Am Ende sollte das nun wie bei mir aussehen. Wobei ich hoffe, dass ihr ein wenig mehr ätzen geübt habt und
eure Skala gleichmäßiger aussieht. Auf dem Foto ist das nicht so gut zu erkennen, live ist die Skala aber sehr
gut zu erkennen und besser abzulesen als die geprägte.

Da wir mit diesem Prozess die Oberfläche des Edelstahl angegriffen haben, sollten wir jetzt nicht einfach
so alles wegräumen. Mindestens solltet ihr den Kessel gut trockenreiben und einen Tag an trockender Luft stehen
lassen. Bei meinem Brauheld hat des gereicht und bisher gab es keinen Rost. Wenn euch eurer Braukessel ein
wenig mehr am Herzen liegt, würde ich ihn komplett passivisieren.

Das geht ganz einfach mit Zitronensäure. Bei den diversen bekannten Online-Händlern und -Portalen kann
einfach Zitronensäurepulver bestellen und daraus eine ca 7%ige Lösung hergestellt werden. Bei einem 30L-Kessel werden
also etwas mehr als 2kg Zitronensäure benötigt. Wurde der Kessel mit der Zitronensäurelösung befüllt,
wird dieser ca. 30 Minuten Zeit zum einwirken gegeben. Imho macht es Sinn in dem Kessel nun andere Gegenstände 
aus Edelstahl (Braupaddel, kleine Schüsseln, Anschlüsse etc.) gleich mit zu passivisieren. Nach den 30 Minuten einfach
die Zitronensäure ablassen (und wenn möglich aufbewahren) und den Kessel und alle Gegenstände gut trocken
reiben und noch ein paar Stunden an der Luft stehen lassen. So sollte euer Edelstahl wieder so robust wie 
zuvor sein.

Natürlich könnt ihr mit dieser Technik alle möglichen Sachen auf Edelstahl ätzen. Wie wäre es mit einem
wunderschönen Warsteiner-Logo auf einem Edelstahl-Growler? Oder einfach der eigene Name, damit es beim
nächsten Hobbybrauerstammtisch keine Verwechslung gibt.
