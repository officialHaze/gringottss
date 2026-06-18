FROM alpine:3

# Set the working directory
WORKDIR /gringottss

# Copy the binary
COPY ./bin/gringottss .

# Copy DB migrations
COPY ./db/migrations ./db/migrations

# Expose the port for this service
# Note:- If the defauly port is changed in settings file, same has to be updated here.
EXPOSE 6230

CMD ["./gringottss", "start_server"]