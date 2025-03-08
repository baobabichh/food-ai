#include <drogon/drogon.h>
#include <drogon/HttpAppFramework.h>
#include <drogon/utils/Utilities.h>
#include <fstream>
#include <iostream>
#include <optional>

#include "nlohmann/json.hpp"

static std::string getFileStr(const std::string& file_path)
{
    std::ifstream file{file_path};
    if (!file.is_open())
    {
        return {};
    }

    std::string file_str{(std::istreambuf_iterator<char>(file)), std::istreambuf_iterator<char>()};
    return file_str;
}

static std::optional<std::string> getStringSafe(const nlohmann::json& object, const std::string& key)
{
    if(object.contains(key) && object[key].is_string())
    {
        return object[key].get<std::string>();
    }
    return std::nullopt;
}


int main()
{
    {
        const auto file_str = getFileStr("/home/goshan9/food-ai/credentials.json");
        if(!nlohmann::json::accept(file_str))
        {
            std::cerr << "Error while config file parse\n";
            return 1;
        }

        nlohmann::json object = nlohmann::json::parse(file_str);
        const auto db_user = getStringSafe(object, "db_user");
        const auto db_pass = getStringSafe(object, "db_pass");

        if(!db_user || !db_pass)
        {
            std::cerr << "Error while config parse, no db_user or db_pass\n";
            return 1;
        }

        drogon::orm::MysqlConfig cfg{};
        cfg.host = "127.0.0.1";
        cfg.port = 3306;
        cfg.databaseName = "food_ai";
        cfg.username = db_user.value();
        cfg.password = db_pass.value();
        cfg.connectionNumber = 1;
        cfg.name = "food_ai";

        drogon::app().addDbClient(cfg);
    }

    drogon::app().addListener("0.0.0.0", 5555);
    drogon::app().loadConfigFile("../config.json");
    drogon::app().run();
    return 0;
}
