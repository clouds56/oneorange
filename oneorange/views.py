from django.shortcuts import render
from django.http import HttpResponse

def home(request):
    return HttpResponse("""<h1>Welcome to one-orange</h1><a href="articles/">articles</a> <a href="admin/">admin</a>""")
