# Zero Effort Hosting Daemon (ZEHD)

## Introducing ZEHD

ZEHD is a static site generator that streamlines the process of converting markdown, org-mode, and gohtml/html files to HTML. Similar to Hugo, ZEHD allows users to automatically create paths and parse their files. The key difference is that ZEHD removes the extra steps required to convert markdown to HTML, as it can do this automatically.

### Key Features

ZEHD has several main features:

1. Building gohtml/html files
2. Converting org-mode files to HTML
3. Converting markdown to HTML
4. Caching to eliminate the need for service restarts
5. Git integration for easy content management

### Building Gohtml/HTML Files

ZEHD supports Go templates with the .gohtml extension, and standard HTML files with the .html extension. These templates can include Go template language constructs such as control structures, loops, and functions, and can be used to generate dynamic content.

### Converting Org-mode Files to HTML

In addition to markdown, ZEHD can also convert org-mode files to HTML. Org-mode is a powerful markup language that supports outlining, notes, lists, tables, and other advanced features.

### Converting Markdown to HTML

Like other static site generators, ZEHD can also convert markdown files to HTML. The conversion is handled by the Blackfriday Markdown processor.

### Caching

One of the key benefits of ZEHD is its caching functionality. This eliminates the need to restart the service when new content is added or updated, and ensures that content is delivered quickly to users.

## Example Templates

Here's an example of a simple Go template:

```html
<!-- layout.gohtml -->
{{ define "layout" }}
<!DOCTYPE html>
<html>
    <head>
        <title>{{template "title"}}</title>
    </head>
    <body>
        {{template "templatePart"}}
    </body>
</html>
{{ end }}

<!-- pagename.gohtml -->
{{define "title"}}
<h1>
    The name of my Title
</h1>
{{end}}

{{define "templatePart"}}
<div>
    <p>
        Whatever HTML you want to have displayed
    </p>
</div>
{{end}}
```

## Kubernetes and Hybrid Setup Ready

ZEHD can be easily configured using environment variables. Here's a comprehensive list of available options:

```bash
BACKEND="http://YourBackend.example.com:8080"
HOSTNAME="ZEHD"
TEMPLATEDIRECTORY="/var/zehd/templates"
TEMPLATETYPE="gohtml"
REFRESHCACHE="60"
PROFILER=false
GITLINK=""
GITUSERNAME=""
GITTOKEN=""
JSPATH=""
CSSPATH=""
IMAGESPATH=""
DOWNLOADSPATH=""
```

### New Git Integration Feature

ZEHD now supports automatic content management through Git integration. By specifying a Git repository URL, ZEHD will automatically clone the repository and fetch any changes you make. This eliminates the need for manual uploads or downloads to update your site content.

To use this feature:

1. Set the `GITLINK` environment variable to your repository URL.
2. If your repository is private, also set `GITUSERNAME` and `GITTOKEN`.
3. ZEHD will automatically clone the repository on startup and periodically fetch updates.

This feature allows for a more streamlined workflow, enabling you to manage your site content directly through your Git repository.

## Docker Deployment

To try or run ZEHD on Docker:

```bash
docker pull zehd/zehd:latest
docker run -d --name zehd \
  -e BACKEND=http://your-backend-url:8080 \
  -e GITLINK=https://github.com/your-username/your-repo.git \
  -e GITUSERNAME=your-username \
  -e GITTOKEN=your-token \
  -p 8080:8080 \
  zehd/zehd:latest
```

Adjust the environment variables as needed for your setup.

## Backend Service

For collecting data on site visitors, check out the backend service: <https://github.com/APoniatowski/zehd-backend>

## Future Plans

- Implementing a worker/agent for dynamic IP/hostname banning
- Adding frontend-to-backend calls to check for banned IPs/hostnames
- Expanding Git integration features
- Inter-service replication, to save time on git pulls/clones or upload/downloads (PLANNED)

## Getting Started

To get started with ZEHD:

1. Clone the ZEHD repository or pull the Docker image
2. Set up your environment variables as described above
3. If using Git integration, ensure your content repository is set up and accessible
4. Run ZEHD
5. Access your site through the configured port (default: 8080)

For more detailed instructions, please refer to the documentation (link to be added).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[https://github.com/APoniatowski/zehd/blob/main/LICENSE]
