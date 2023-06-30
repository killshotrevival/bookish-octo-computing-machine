FROM --platform=linux/amd64 golang:1.19.5 as build
WORKDIR /temp
COPY . /temp/

RUN CGO_ENABLED=0 go build -v

FROM --platform=linux/amd64 registry.digitalocean.com/getastra/furious:0.0.1
COPY --from=build /temp/endgame /endgame
CMD [ "/endgame" ]