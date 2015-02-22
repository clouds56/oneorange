from django.conf.urls import url

from articles import views

urlpatterns = [
    url(r'^$', views.index, name = 'index'),
    url(r'^(?P<author>[^/]+)/?$', views.author, name='author'),
    url(r'^(?P<author>[^/]+)/(?P<anthology>[^/]+)/?$', views.anthology, name='anthology'),
    url(r'^(?P<author>[^/]+)/(?P<anthology>[^/]+)/(?P<article>\d+)/?$', views.article, name='article'),
]
