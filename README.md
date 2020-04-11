# Search Engine

An implementation of a search engine written in golang.
Currently only contains an indexer, ranker with a basic parser (tokenizer) and a "crawler".
Documents are indexed in an inverted index and k-gram index for different query methods.

Crawler uses the [MediaWiki action API](https://www.mediawiki.org/wiki/API:Main_page) to scrape the introductory paragraph, 
and uses out-going links in the article to find more Wikipedia articles.

### Reference

Christopher D. Manning, Prabhakar Raghavan, and Hinrich Schütze. 2008. Introduction to Information Retrieval. Cambridge University Press, USA.