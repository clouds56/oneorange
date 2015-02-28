from django.shortcuts import render, redirect
from django.http import HttpResponse
from django.template import RequestContext, loader
from django.contrib import auth
from django.contrib.auth.models import User
from django.contrib.auth.decorators import login_required

from common.util import get_referer_view, get_relative_url
from articles.models import Author, Article, Anthology
# Create your views here.

def index(request):
    authors = Author.objects.all()
    template = loader.get_template("base.html")
    context = RequestContext(request, {'authors': authors})
    return HttpResponse(template.render(context))

def login(request):
    referer = get_referer_view(request, "/articles")
    if referer.find("/accounts/login") == 0:
        referer = "/articles"

    if request.method == 'GET':
        if not 'next' in request.GET or not request.GET['next']:
            return redirect("/accounts/login?next="+referer)
        return render(request, "login.html", {'next': request.GET['next']})

    next = ""
    if 'next' in request.POST and request.POST['next']:
        next = get_relative_url(request.POST['next'], request.META.get('SERVER_NAME'), referer)
    elif 'next' in request.GET and request.GET['next']:
        next = get_relative_url(request.GET['next'], request.META.get('SERVER_NAME'), referer)
    else:
        next = referer
    if next.find("/accounts/login") == 0:
        next = referer;llllkkk
    if request.user.is_authenticated():
        return redirect(next)
    if not 'username' in request.POST:
        return render(request, "login.html", {'error_msg': "No username", 'next': next})
    if not 'password' in request.POST:
        return render(request, "login.html", {'error_msg': "No password", 'next': next})
    username = request.POST['username']
    if User.objects.filter(username=username).exists():
        return render(request, "login.html", {'error_msg': "User not exist", 'next': next})
    password = request.POST['password']
    user = auth.authenticate(username=username, password=password)
    if not user:
        return render(request, "login.html", {'error_msg': "User or password wrong", 'next': next})
    auth.login(request, user)
    return redirect(next)

@login_required
def logout(request):
    auth.logout(request)
    return redirect(request.META.get('HTTP_REFERER'))

def signup(request):
    if not 'csrfmiddlewaretoken' in request.POST:
        return render(request, "signup.html", {'error_msg': "New register"})
    if not 'username' in request.POST:
        return render(request, "signup.html", {'error_msg': "No username"})
    if not 'password' in request.POST:
        return render(request, "signup.html", {'error_msg': "No password"})
    if not 'email' in request.POST:
        return render(request, "signup.html", {'error_msg': "No email"})
    username = request.POST['username']
    email = request.POST['email']
    password = request.POST['password']
    if User.objects.filter(username=username).exists():
        return render(request, "signup.html", {'error_msg': "Username exists", 'last': {'username': username, 'email':email}})
    if User.objects.filter(email=email).exists():
        return render(request, "signup.html", {'error_msg': "Email exists", 'last': {'username': username, 'email':email}})
    user = User.objects.create_user(username, email, password)
    user.save()
    return redirect("/accounts/step2/"+username)

def step2(request, author):
    if not request.user.is_authenticated():
        return redirect("/accounts/login")
    if not 'csrfmiddlewaretoken' in request.POST:
        return render(request, "step2.html", {'author': author})
    username = request.user.username
    if not request.POST['author']:
        return render(request, "step.html", {'error_msg': "No author name", 'author': author})
    author = request.POST['author']
    if Author.objects.filter(name=author).exists():
        return render(request, "step.html", {'error_msg': "Author name exists", 'author': author})
    anauthor = Author.objects.create(name=author, user=request.user)
    anauthor.save()
    return redirect("/accounts/success")
