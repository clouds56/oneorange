from django.conf.urls import patterns, include, url
from django.contrib import admin
from django.views.generic import TemplateView

admin.autodiscover()

urlpatterns = patterns('',
    # Home Page -- Replace as you prefer
    #url(r'^$', TemplateView.as_view(template_name='home.html'), name='home'),
    url(r'^$', 'oneorange.views.home', name='home'),

    url(r'^articles/', include('articles.urls')),
    url(r'^accounts/', include('accounts.urls')),
    #url(r'^admin/doc/', include('django.contrib.admindocs.urls')),
    url(r'^admin/', include(admin.site.urls)),
)
