FROM alpine:3

# Set the working directory
WORKDIR /gringottss

# Copy the binary
COPY ./bin/gringottss .

# Copy DB migrations
COPY ./db/migrations ./db/migrations

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Give USER and GROUP appropriate perms
RUN mkdir -p /var/log && chown -R appuser:appgroup /var/log
RUN chown -R appuser:appgroup /gringottss
RUN chown -R appuser:appgroup /home

# Expose the port for this service
# Note:- If the defauly port is changed in settings file, same has to be updated here.
EXPOSE 6230

# Switch to non-root user
USER appuser

CMD ["./gringottss"]