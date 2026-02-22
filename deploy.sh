#!/bin/bash
rsync -av --delete --exclude 'app.db' dist/ admin@runeco.de:~/apps/social.runeco.de/
