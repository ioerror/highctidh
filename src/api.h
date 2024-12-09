#ifdef _WIN32
#    ifdef DLL_EXPORT
#        define HIGHCTIDH_API __declspec(dllexport)
#    else
#        define HIGHCTIDH_API __declspec(dllimport)
#    endif
#else
#    define HIGHCTIDH_API
#endif
