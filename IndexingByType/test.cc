#include <iostream>
#include <string>
#include "data_whale.h"

using std::cout;

class Foo {
  public:
    Foo(int num) : num_(num) {}
    void Print() {
      cout << num_ << "\n";
    }

  private:
    int num_;
};

class Bar {
  public:
    Bar(const char* text) : text_(text) {}
    void Print() {
      cout << text_ << "\n";
    }

  private:
    std::string text_;
};


int main() {
  DataWhale::Register(new Foo(100));
  DataWhale::Register(new Bar("hello world"));
  DataWhale::Get<Foo>()->Print();
  DataWhale::Get<Bar>()->Print();
}
