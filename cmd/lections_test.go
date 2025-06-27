package main

import (
	"testing"
)

func TestSplitByPunctuationAndLength(t *testing.T) {
	rawHtml := `
	<p>Hermanos: Llevamos un tesoro en vasijas de barro, para que se vea que esta fuerza tan extraordinaria proviene de Dios y no de nosotros mismos. Por eso sufrimos toda clase de pruebas, pero no nos angustiamos. Nos abruman las preocupaciones, pero no nos desesperamos. Nos vemos perseguidos, pero no desamparados; derribados, pero no vencidos.<br>
	<br>
	Llevamos siempre y por todas partes la muerte de Jesús en nuestro cuerpo, para que en este mismo cuerpo se manifieste también la vida de Jesús. Nuestra vida es un continuo estar expuestos a la muerte por causa de Jesús, para que también la vida de Jesús se manifieste en nuestra carne mortal. De modo que la muerte actúa en nosotros, y en ustedes, la vida.<br>
	<br>
	Y como poseemos el mismo espíritu de fe que se expresa en aquel texto de la Escritura: <em>Creo, por eso hablo,</em> también nosotros creemos y por eso hablamos, sabiendo que aquel que resucitó a Jesús nos resucitará también a nosotros con Jesús y nos colocará a su lado con ustedes. Y todo esto es para bien de ustedes, de manera que, al extenderse la gracia a más y más personas, se multiplique la acción de gracias para gloria de Dios.</p>`

	chunks := splitByPunctuationAndLength(rawHtml)

	if len(chunks) < 2 {
		t.Errorf("Expected multiple chunks for long reading, but got only %d chunk(s)", len(chunks))
	}

	for i, chunk := range chunks {
		if len(chunk) > 225 {
			t.Errorf("Chunk %d is too long (%d characters): %s", i+1, len(chunk), chunk)
		}
	}

	for i, c := range chunks {
		t.Logf("Chunk %d:\n%s\n", i+1, c)
	}
}

