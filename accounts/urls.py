from django.conf.urls import url

from accounts import views

urlpatterns = [
    #url(r'^$', views.index, name = 'index'),
    url(r'^login/?$', views.login, name = 'login'),
    url(r'^logout/?$', views.logout, name = 'logout'),
    url(r'^signup/?$', views.signup, name = 'signup'),
    url(r'^step2(/(?P<author>[^/]+))?/?$', views.step2, name = 'step2'),
    url(r'^success/?$', views.success, name = 'success'),
]
