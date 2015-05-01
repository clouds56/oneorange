from django.conf.urls import url

from articles import views

urlpatterns = [
    url(r'^$', views.index, name = 'index'),
    url(r'^(?P<author_name>[^/]+)/?$', views.AuthorView.as_view(), name='author'),
    url(r'^(?P<author_name>[^/]+)/(?P<anthology_name>[^/]+)/?$', views.AnthologyView.as_view(), name='anthology'),
    url(r'^(?P<author_name>[^/]+)/(?P<anthology_name>[^/]+)/new/?$', views.newarticle, name='anthology'),
    url(r'^(?P<author_name>[^/]+)/(?P<anthology_name>[^/]+)/(?P<article_id>[^/]+)/?$', views.article, name='article'),
]
