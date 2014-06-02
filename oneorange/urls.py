from django.conf.urls import patterns, include, url

from django.contrib import admin
admin.autodiscover()

urlpatterns = patterns('',
    # Examples:
    url(r'^$', 'oneorange.views.home', name='home'),
    url(r'^comments/', include('django.contrib.comments.urls')),
    # url(r'^blog/', include('blog.urls')),

    url(r'^articles/', include('articles.urls')),
    url(r'^admin/', include(admin.site.urls)),
)
