package rawio

import (
"bytes"
	"bufio"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"fmt"
)

func init() {
	// Ignore EOF so we get the same
	// 'stream' behaviour like NNTP
	//testNoEOF = true
}

type report interface {
	Error(args ...interface{})
}

// Cancel upload
func TestEmpty(t *testing.T) {
	r := New(bufio.NewReader(strings.NewReader(".\r\n")), true)
	_, e := io.Copy(ioutil.Discard, r)
	if e != nil {
		t.Error(e)
	}
}

// Wrong input but with dot
func TestSimple(t *testing.T) {
	raw := "Blabla\r\nMore blablabla.\r\nBlaaaaaaaat"
	raw += "\r\n.\r\n"

	r := New(bufio.NewReader(strings.NewReader(raw)), false)
	n, e := io.Copy(ioutil.Discard, r)
	if e != nil {
		t.Error(e)
	}
	if n != int64(len(raw)) {
		t.Errorf("Bytes mismatch, expect=%d received=%d", len(raw), n)
	}
}

func TestSimpleStreamMODE(t *testing.T) {
	raw := "Blabla\r\nMore blablabla.\r\nBlaaaaaaaat"
	raw += "\r\n.\r\n"
	l := len(raw)

	// STREAM-mode
	raw += "Blabla\r\nMore blablabla.\r\nBlaaaaaaaat"
	raw += "\r\n.\r\n"	

	inbuf := bufio.NewReader(strings.NewReader(raw))

	b := new(bytes.Buffer)
	r := New(inbuf, false)
	n, e := io.Copy(b, r)
	if e != nil {
		t.Error(e)
	}
	if n != int64(l) {
		fmt.Printf("%s\n", b)
		t.Errorf("Bytes mismatch, expect=%d received=%d", len(raw), n)
	}

	// Second
	b = new(bytes.Buffer)
	r = New(inbuf, false)
	n, e = io.Copy(b, r)
	if e != nil {
		t.Error(e)
	}
	if n != int64(l) {
		fmt.Printf("%s\n", b)
		t.Errorf("Bytes mismatch, expect=%d received=%d", len(raw), n)
	}
}

// Correct server reply
func TestArticle(t *testing.T) {
	raw := strings.Replace(`Path: asa019.ams.xsnews.nl!feeder04.ams.xsnews.nl!border04.ams.xsnews.nl!feed.xsnews.nl!fbe002.ams.xsnews.nl!newsfeed.fsmpi.rwth-aachen.de!newsfeed.straub-nv.de!feeder.erje.net!1.eu.feeder.erje.net!feeder2.ecngs.de!ecngs!feeder.ecngs.de!81.171.118.63.MISMATCH!peer03.fr7!news.highwinds-media.com!peer01.am1!peering.am1!npeersf04.am4!fx15.fr7.POSTED!not-for-mail
From: Benny15 <viwLYb9mUsLfXjBWLfaUmT8zBdetKvSsByfnOqwNDcwZHC-sU-pJr3s-pXaKnM1QoJ3.EWuMVOFvoWJvlkJGD78L115aqhXdV1GNSjBIKeB84pWpCRNIjolrov5InLJB9-sXy@17a10b09c12d13z00.6907067771.10.1429970067.1.NL.TVOnL6991U0mmb17M0eDJ3qyF7lr2gYwm53mULq-s2fEVsChAlh1RP8aPxdcIS6z-p>
Subject: Das gro?e Hansi Hinterseer Open Air 2011 | Jozefien
Newsgroups: free.pt
Message-ID: <lUdFcpPeOBskJw7VQOO7E@spot.net>
X-XML: <Spotnet><Posting><Key>7</Key><Created>1429970067</Created><Poster>Benny15</Poster><Tag>Jozefien</Tag><Title>Das gro&#223;e Hansi Hinterseer Open Air 2011</Title><Description>zelf toen opgenomen van de TV[br][br]10.09.2011 [br]20:15 Uhr Das gro&#223;e Hansi Hinterseer Open Air 2011 Aufzeichnung aus Kitzb&#252;hel | Das Erste Tipp [br]Wieder begr&#252;&#223;t Publikumsliebling Hansi namhafte Musikkollegen im Tennisstadion seiner Heimatstadt Kitzb&#252;hel. Unter anderem sind Andrea Berg, die Zellberg Buam, Encho Keryazow, Rob Spencer, das Original Tiroler Echo und das Alpentrio Tirol dabei.[br] [br]Details Wiederholungen [br]Mitrei&#223;ende Stimmungshits, ber&#252;hmte Evergreens und traditionelle Tiroler Lieder stehen f&#252;r einen gelungenen Abend. Auch der Gastgeber selbst wird einige seiner Hits singen.[br]Zwei Tage vor seinem gro&#223;en Open Air wanderte Hansi Hinterseer auch 2011 wieder mit
X-XML:  seinen Fans bei viel Musik und guter Laune auf den Kitzb&#252;heler Hausberg &quot;Hahnenkamm&quot;. Dass er bei Menschen aller Altersklassen gleicherma&#223;en beliebt ist, bewiesen wieder die Tausende von Fans, die sich diesen Event nicht entgehen lassen wollten. Auch die H&#246;hepunkte der Fanwanderung zeigt Das Erste im Rahmen dieses &quot;Gro&#223;en Hansi Hinterseer-Open Air 2011&quot;.[br][br]dolby digital[br][br]speelduur +- 2 uur</Description><Website>http://www.tvmovie.de/das-grosse-hansi-hinterseer-open-air-2011-2150274.html?image-number=2</Website><Image Width='1280' Height='720'><Segment>mJOBBkxZsH8kpw7VQ1UCs@spot.net</Segment></Image><Size>6907067771</Size><Category>01<Sub>01a10</Sub><Sub>01b09</Sub><Sub>01c12</Sub><Sub>01d13</Sub><Sub>01z00</Sub></Category><NZB><Segment>VlqUl8F9aTQkZw7VQUHNs@spot.net</Segment></NZB></Posting></Spotnet>
X-XML-Signature: fk5KOyiTdcETzguaoh5Khqp8w98JWQ1yJi6jHn3-sLT5Ys5ETnSA8L9CP-s8b8FNLs
X-User-Key: <RSAKeyValue><Modulus>viwLYb9mUsLfXjBWLfaUmT8zBdetKvSsByfnOqwNDcwZHC/U+Jr3s+XaKnM1QoJ3</Modulus><Exponent>AQAB</Exponent></RSAKeyValue>
X-User-Signature: EWuMVOFvoWJvlkJGD78L115aqhXdV1GNSjBIKeB84pWpCRNIjolrov5InLJB9-sXy
X-User-Avatar: iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAAEzkAABM5AY/CVgEAAAstSURBVFhHrZcHUFRp9sXZ3UHANKNIhoYGoQGboCgSFFsJgoKYcBABCaKCiKiIoghiwAFRUVgVxATqDAMGDCOiLIZRZwyjmAnaEo3orLqrW1vTv/20+l/zr9qtDdacqq5Xr/rre86999z7Xmv8O6hUqm5A787OF3Y///zOStxbl+7dnX1s38Yn9Sd2nBbf/V599LdHff31NZfPnVCerSp+emzvWo6Xrfnb/k2L3ucsi6ByRQSnijMRgizVx387vFWpDO9cO3f6u91r+HptHAfXzmR7ZiTpcX4sCB/Gj/uWsjctgqKshA8CrNU/+20gSvqHcycPHT2yaSHFiYFsnR9CQeIE8hImEhPkTnqCP5d3pZEw0Yt1S6J/EQKs1D/9dNy9cSXizImKHbWVJT/eqau8X7U5jfmTPckI82bN9NGkhytYFjaCeZM92J46ifPbU6jemUnF1mW/wDupOsyn40x58dHmUwWcP5jL7sxphI2yJ8zHgbhxQ0SmnsQEODM/1IOkUC8K5o3j27w4zu1I5cdDuTQ2/zRVHebTUZ6/alVxcjB3qtdz7cBKMmO9WRjixozxbiyOHMWMQBeWif4nR/qwJnE8+1dGcSovlrKsUE4d2XZLHebTcb66cs6qSF/2JI/leG4UpUsnsjEhkFUzA5gTPopNC0IoXBTChoWh7F0Zy4HsKI7lTKc4NYgDJUs/GDFIHerToOxQykqz4n+Z5WbOlgQftgghW1KC2b8inLiQEexZMYMDq2PZtmQKFRvmECHaUpgcSGXWVMpzY6gqy/7rD7X7j79
X-User-Avatar: //9ZfHfJ/x9dbc04ED5GSN3MUexeNoyJjMnX5s8hdEMZmMXKHc+LYlR5G4dKpFK6IYU9GOFsWBvNtfjzzovyYLTyyZ/Mirn9feelvKtUAddj/HpcvnfffkBpJdpSCpGAXIvwcWSTcX5oVTf7cIKqyYyhaMJ6qknSOiAWUET+RCf7OFGSEEfOlghkT3UgIcaUyJ5rj5VtevnzzxkEd+r+DmP/PdhVktYYoZPgPtSTER87K+DFcrFhP0ZJp5EYNJzfRn32FqZSsSyElbBRpUd7ETPNl8lhXlseM5qvEsewW1cuIH83lE9ufC2/YqMP/Z7x7985q1swv35lKeuPpZsPy5BCyk6dQUpBJYXYKBQm+bE70oyh7DuuXRLNQCJgvjDsmwI1FogVFqSHsT5tIfnIQy2JGsjsjhMrSnJdvnim/VFP8e6hU7+2WLE58qVA4ULM3i7YbVRzZm0fBilhy0mLYNT+AP2VNZuu8iRzcmMKS2CB8/IYyPzZQVCicdUnBFMwNYOP88aTP8CE3IYAFkSP4Wpy9UlWeqfr56ejXoKum+9e4d0854tLpind1pUu5dCCPSwe3sCEznqWhQ9kVN4wf1k2nXDwTqvdkM2WiN/OmB1CwOJKD+cnkL49mcbQvGdEK1sb7iasPc8XiSgrzYEdMEMqIqTTGxzzrLM2vuddy00VN+c84tG/zobz0UHbnJVK0ahZzpvtRvGYWG5PDSBxmx/Fv8ilev5igUc6UpkdSmhnNzuWRzI4Lwl0hZ7S3PbOneDJhzGDmTXFnzoQhzBnnSomjnJvewyE1hsY9Xx1W0/0z/jQzaU++jS0xwxwY0d8Qh/792LYunu+rdxMX6suu/FQWxo5j69JwSlbPZNvKWaycHcxaDycKJ3ixOzqAo0mTODkvlOvJoTQsj6
X-User-Avatar: N1UQydUePpnOTLk4FyWpfOfvtapeqnpvwV91Sqnle9grvqNPQo/EM/cnvpM99An/HOUsJEFqnTAxnhasvqhCDufbeBs0f/yIoIf0qNpDxxc+XPIsNXvl50BXrxcsIoXkzyoSs+hOcLIukMGEGbqxMdcnteTR1Ha+e14WraX3Hrm0rFXUs3Lv/OiLOaplT1NafM0ILlxkZM7N6TcZo6TPdzorZkMbdqt1FTkUuGpQ23dSQ8kjvwYKgLSmcHWoYOpMVrMG3eQ+mY7E1HSjidk/1o8xhE50Bn/uI9DGXZhpVq2l9xY+6Sooc9bbjQzZTvtIw5aWxB+RcSdpqZU9TdiJxeRsJoCr7/JoOmi2X8cUkEuX1MaNGT0SKxo3mwIw9Fr1ucHWl1G0zrcCFirCfPU8N4FTue58MG82SwMy/s7Hi4IvnBbeihptbQaFWptG95h7Q2aZhwsJsRK/X0qNWVkNtTn7Vf6LNTCCrWNGRp4BB2romkcstCIl3tOdnThA5je1pM7Hlk60C7g5znTk4iUydaXAfSoHChNkRBRbCCnY623Bwg54mVjDfxEbT9/flINb2Gxu3qU6Nb5CO5+3sTyrWN2d7LgDwdfWb37UuOEJDSXZeiboak9jMm2NeRCd5yUnrqcVqU/7phf9rNHeiUOXLGWkappZRiuQ3lgxw562DHntHuXE6P53yQDxVSczrtHXjvp6C1uiRRTS9eQJPT8h/3seN2N4kQYMgyLV2W99GnTE/Cju6mpHfXZ1d3YzZpGZLUrS+FmgZc1bFgX28zTupLeWwup91WzjHRrroPpZ8ylo4Ab24NdOCUlwsvk6Zw0WsQtbL+tMkdeSvOPsxIqvxIfkU8B+4FhTW1appRr21OeXcDDhlIWCfKXyJKv0vbhBU6uuwXlSnrbsIxc
X-User-Avatar: d/cTUpTD0uqe0s4Z2T1sQIdVnIuWlhTZSGlxsqC0zIr7rs4c9DRhuIxHpR4u3FcJuOpbABPDC3pjA9/26lSSTQeXbjg/mSIH0pRznodc0q66XK0lzE5Iut8bQP2CsLVogVHhJhSHWO+E0IatKU0f96fk33MOKNrTrvJANqlcn40teSwGLUjY72oU3hwXz6AGtH7G9HBKGNDKJda0iyx4bGhNW8DfOm4WxOq0ZixJve1mTOPelryU28LykwsqO1mwjYdPVYKMZVCQJ52X04JAeU6JpzTltDeU5RSV8bpzyWc7W3OSxM5Ly0cuWQspVb44LWnO88GDhKmdOKoqEj9hx3g785hGytumst4YmLLGxth2pKvyjUeTIlpfvG5BcrPrbjbS8pVI1sufGbGYV0zNolg34qMC3vrcVbblAphxN3CnAd6GHC5nxXnvjCnVFtP7AwJtQaW3DC2YV8/M8rMzNhhLeWSzI6LwgdXhzrTOdKFuyNdUfaX0ybOPdW1oH1O1GMN5bAxPBdG+yCgsYeUm59J+EH4oVbHjBN61tT0sOCMlT2Xtcy4Yu1AnYcnlR6DqRZL6Ja5LbUi26rhntS4uHLH0pZrYtFc9/OmwVdBs9yJdmcn2sRYtrmL8RzkwGMLW9oMbXgm2tkxbiwaTXn5+S+H+dPVy5wmTRNuCR9c1ZLwvdiGdZ8Zc1HLgp80JdwWBn2ga0OrqTBRf0dhvAG0iM9jh6E8s3cV5hKbTmr/cRyfWooz5nZi5sW9pR0d/e1oF0uqQ+yJTiNLXulaCg8E0bJjy9qPk9BQVxuqTFhQ/3DgCDr0bLkv9sFPGoZcFNcfBPk1LXOxcs1p7m3FA2E+pbh+8ECrvtiCusIPom3tpna0mQkiiS0dFnaCeAAdEhmPJdY8NbWiy8KG1y5D
X-User-Avatar: ee4fSGdaWv3D+ouRH8n/D+nwu1tnTivuL1+99d7U2Y+uDvJ51yRT0KQ/iLs6Mho0rVBqCQGaFjzSNEfZw4qWXta0iHFsE+3r6COjs48NXQZ2vBDb8c/Sgbx29OSVZwCPxkx4/WjO3FtNeXn5DXfqh3/gUtP+ayhBq0H87Wo8cDjs6rKstMbV6083JKZduDZ22pXrQeFdN3y/7GpQTOpqHB7c1ege1NUwclLX7cCpXXfGhz+9GTztUkf66gsPVmUfuL9ufXxTTe3kBpXKWLx4qqP/f2ho/APt2Mj67AOBQgAAAABJRU5ErkJggg
Content-Type: text/plain; charset=ISO-8859-1
Content-Transfer-Encoding: 8bit
Lines: 13
X-Complaints-To: http://www.newsleecher.com/support/
NNTP-Posting-Date: Sat, 25 Apr 2015 13:54:26 UTC
Date: Sat, 25 Apr 2015 13:54:26 GMT
X-Received-Body-CRC: 2876841174
X-Received-Bytes: 8002

zelf toen opgenomen van de TV

10.09.2011 
20:15 Uhr Das gro?e Hansi Hinterseer Open Air 2011 Aufzeichnung aus Kitzb?hel | Das Erste Tipp 
Wieder begr??t Publikumsliebling Hansi namhafte Musikkollegen im Tennisstadion seiner Heimatstadt Kitzb?hel. Unter anderem sind Andrea Berg, die Zellberg Buam, Encho Keryazow, Rob Spencer, das Original Tiroler Echo und das Alpentrio Tirol dabei.
 
Details Wiederholungen 
Mitrei?ende Stimmungshits, ber?hmte Evergreens und traditionelle Tiroler Lieder stehen f?r einen gelungenen Abend. Auch der Gastgeber selbst wird einige seiner Hits singen.
Zwei Tage vor seinem gro?en Open Air wanderte Hansi Hinterseer auch 2011 wieder mit seinen Fans bei viel Musik und guter Laune auf den Kitzb?heler Hausberg "Hahnenkamm". Dass er bei Menschen aller Altersklassen gleicherma?en beliebt ist, bewiesen wieder die Tausende von Fans, die sich diesen Event nicht entgehen lassen wollten. Auch die H?hepunkte der Fanwanderung zeigt Das Erste im Rahmen dieses "Gro?en Hansi Hinterseer-Open Air 2011".

dolby digital

speelduur +- 2 uur.`, "\n", "\r\n", -1)
	raw += "\r\n.\r\n"
	r := New(bufio.NewReader(strings.NewReader(raw)), false)
	n, e := io.Copy(ioutil.Discard, r)
	if e != nil {
		t.Error(e)
	}
	if n != int64(len(raw)) {
		t.Errorf("Bytes mismatch, expect=%d received=%d", len(raw), n)
	}
}