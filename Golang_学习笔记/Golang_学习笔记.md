# 1、基本语法

## 1.1、基本数据类型

- bool

- string

- int  int8  int16  int32  int64
  uint uint8 uint16 uint32 uint64 uintptr

- byte // uint8 的别名

- rune // int32 的别名
      // 表示一个 Unicode 码点

- float32（少用）  float64（常用）

- complex64  complex128

`int`, `uint` 和 `uintptr` 在 32 位系统上通常为 32 位宽，在 64 位系统上则为 64 位宽。 当需要一个整数值时应使用 `int` 类型，除非你有特殊的理由使用固定大小或无符号的整数类型。















