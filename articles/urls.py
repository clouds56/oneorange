from django.conf.urls import url

from articles import views

urlpatterns = [
    url(r'^$', views.index, name = 'index'),
    url(r'^(?P<author>[^/]+)/?$', views.archive, name='archive'),
    url(r'^(?P<author>[^/]+)/(?P<id>\d+)/?$', views.detail, name='detail'),
]
