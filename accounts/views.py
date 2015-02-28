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
        next = get_relative_url(request.GET['next'], request.META.get('SERVER_NAME'), referer)
        if next.find("/accounts/login") == 0:
            next = referer
        return render(request, "login.html", {'next': next})

    if 'next' in request.POST and request.POST['next']:
        next = get_relative_url(request.POST['next'], request.META.get('SERVER_NAME'), referer)
    elif 'next' in request.GET and request.GET['next']:
        next = get_relative_url(request.GET['next'], request.META.get('SERVER_NAME'), referer)
    else:
        next = referer
    if next.find("/accounts/login") == 0:
        next = referer
    if request.user.is_authenticated():
        return redirect(next)
    if not 'username' in request.POST:
        return render(request, "login.html", {'error_msg': "No username", 'next': next})
    if not 'password' in request.POST:
        return render(request, "login.html", {'error_msg': "No password", 'next': next})
    username = request.POST['username']
    if not User.objects.filter(username=username).exists():
        return render(request, "login.html", {'error_msg': "User not exists", 'next': next})
    password = request.POST['password']
    user = auth.authenticate(username=username, password=password)
    if not user:
        return render(request, "login.html", {'error_msg': "Wrong username or password", 'next': next})
    elif not user.is_active:
        return render(request, "login.html", {'error_msg': "User is deleted", 'next': next})
    auth.login(request, user)
    return redirect(next)

@login_required
def logout(request):
    auth.logout(request)
    referer = get_referer_view(request, "/articles")
    if referer.find("/accounts/logout") == 0:
        referer = "/articles"
    return redirect(referer)

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
    user = auth.authenticate(username=username, password=password)
    auth.login(request, user)
    return redirect("/accounts/step2/"+username)

@login_required
def step2(request, author):
    username = request.user.username
    if request.method == 'GET':
        return render(request, "step2.html", {'username': username, 'author': author})
    if not request.POST['author']:
        return render(request, "step2.html", {'error_msg': "No author name", 'username': username, 'author': author})
    author = request.POST['author']
    if Author.objects.filter(user__username=username).exists():
        return render(request, "step2.html", {'error_msg': "Author already created", 'username': username, 'author': author})
    if Author.objects.filter(name=author).exists():
        return render(request, "step2.html", {'error_msg': "Author name exists", 'username': username, 'author': author})
    anauthor = Author.objects.create(name=author, user=request.user)
    anauthor.save()
    return redirect("/accounts/success")
