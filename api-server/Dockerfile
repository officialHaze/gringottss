FROM fedora:42

# Install system dependencies, Python, and Node.js via DNF
RUN dnf update -y && dnf install -y --skip-unavailable

# Set the working directory
WORKDIR /gringottss

# Copy the binary
COPY ./bin/gringottss .

# Copy settings
#COPY ./settings/settings.conf.yaml ./settings/settings.conf.yaml

# Copy CONTRIBUTORS
#COPY ./CONTRIBUTORS.md .

# Copy AUTHOR
COPY ./AUTHOR.md .

# Copy DB migrations
COPY ./db/migrations ./db/migrations

RUN groupadd -r appgroup && useradd -r -g appgroup appuser

# Give USER and GROUP appropriate perms
RUN mkdir -p /var/log && chown -R appuser:appgroup /var/log
RUN chown -R appuser:appgroup /gringottss
RUN chown -R appuser:appgroup /home

# Expose the port for this service
EXPOSE 6230

# Switch to non-root user
USER appuser

CMD ["./gringottss"]