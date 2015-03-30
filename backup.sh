#!/bin/sh

python manage.py dumpdata --indent 2 auth.user articles > backup.json
