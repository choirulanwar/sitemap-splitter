# Sitemap Splitter

sitemap-splitter is a lightweight Go package that helps break down large XML sitemaps into smaller, more manageable chunks while maintaining SEO best practices. It automatically creates a sitemap index file that points to the split sitemap files.

Key Features:

- Splits large sitemaps based on a configurable URL limit
- Supports both absolute and relative file paths
- Preserves all URL attributes (lastmod, changefreq, priority)
- Automatically generates a sitemap index file
- Follows sitemap protocol specifications

Example use cases:

- Breaking down sitemaps that exceed the 50MB/50,000 URL limit
- Improving sitemap management for large websites
- Optimizing sitemap loading and processing
