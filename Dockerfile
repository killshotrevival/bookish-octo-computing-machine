FROM --platform=linux/amd64 golang:1.19.5 as build
WORKDIR /temp
COPY . /temp/
RUN CGO_ENABLED=0 go build -v

#####################################################################################################################

FROM --platform=linux/amd64 registry.digitalocean.com/getastra/furious:0.0.1
RUN go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@v2.6.0

COPY --from=build /temp/endgame /endgame
COPY ./resources/subdomain-takeover/fingerprints.json /fingerprints.json
CMD [ "/endgame" ]