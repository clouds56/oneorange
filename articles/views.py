from django.shortcuts import render
from django.http import HttpResponse
from django.template import RequestContext, loader

from articles.models import Article
# Create your views here.

def index(request):
    articles = Article.objects.all()
    template = loader.get_template("archive.html")
    context = RequestContext(request, {'articles': articles})
    return HttpResponse(template.render(context))
