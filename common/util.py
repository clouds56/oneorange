import re

def get_referer_view(request, default=None):
    ''' 
    Return the referer view of the current request

    Example:

        def some_view(request):
            ...
            referer_view = get_referer_view(request)
            return HttpResponseRedirect(referer_view, '/accounts/login/')
    '''

    # if the user typed the url directly in the browser's address bar
    referer = request.META.get('HTTP_REFERER')
    if not referer:
        return default

    # remove the protocol and split the url at the slashes
    return get_relative_url(referer, request.META.get('SERVER_NAME'), default)

def get_relative_url(url, server=None, default=None):
    if not url:
        return default

    # remove the protocol and split the url at the slashes
    if not re.match(url, 'https?:\/\/'):
        return url

    urltree = re.sub('^https?:\/\/', '', url).split('/')
    if server and urltree[0] != server:
        return url

    # add the slash at the relative path's view and finished
    url = u'/' + u'/'.join(urltree[1:])
    return url