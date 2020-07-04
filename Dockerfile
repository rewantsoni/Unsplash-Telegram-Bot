FROM alpine

# Install the required packages
RUN apk add --update git go musl-dev
# Install the required dependencies
RUN go get github.com/go-telegram-bot-api/telegram-bot-api
# Setup the proper workdir
WORKDIR /root/Unsplash
# Copy indivisual files at the end to leverage caching
COPY ./LICENSE ./
COPY ./main.go ./
COPY ./utils.go ./
COPY ./types.go ./
RUN go build

#Executable command needs to be static
CMD ["/root/Unsplash/Unsplash"]