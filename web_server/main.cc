#include <drogon/drogon.h>
int main()
{
    drogon::app().addListener("0.0.0.0", 5555);
    drogon::app().loadConfigFile("../config.json");
    drogon::app().run();
    return 0;
}
