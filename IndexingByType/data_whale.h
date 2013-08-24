#include <map>

using std::map;

class DataWhale {
  public:
    template<typename T>
      static T* Get() {
        int type_num = TypeIndexer<T>::type_num;
        return reinterpret_cast<T*>(data_[type_num]);
      }

    template<typename T>
      static void Register(T* datum) {
        int type_num = TypeIndexer<T>::type_num;
        data_[type_num] = (void *)datum;
      }

  private:
    template<typename T>
      struct TypeIndexer {
        static int type_num;
      };

    static int begin_num_;
    static map<int, void*> data_;
};

map<int, void*> DataWhale::data_;

int DataWhale::begin_num_ = 0;

template<typename T>
int DataWhale::TypeIndexer<T>::type_num =  ++DataWhale::begin_num_;
