.PHONY: deps site clean

.DEFAULT_GOAL := site

deps:
	cargo install --git https://github.com/geoffjay/okf --branch okf-web okf

# Build the knowledge base into a static HTML site (output: docs/site/).
site:
	okf site docs/knowledgebase

# Remove the generated site.
clean:
	rm -rf docs/knowledgebase/site
