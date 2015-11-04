all: build

build:	# Run pdflatex twice to resolve references and stuff.
	pdflatex GoAbstractSyntaxTreeVisualizer.tex && biber GoAbstractSyntaxTreeVisualizer && pdflatex GoAbstractSyntaxTreeVisualizer.tex && pdflatex GoAbstractSyntaxTreeVisualizer.tex

clean:
	rm *.aux *.toc *.log *.bcf *.run.xml *.lof *.lot
