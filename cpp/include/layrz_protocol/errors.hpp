#pragma once
#ifndef __LAYRZ_PROTOCOL_ERRORS_HPP__
#define __LAYRZ_PROTOCOL_ERRORS_HPP__

#include <string>
#include <utility>

namespace layrz::protocol {

enum class Error {
    Ok,
    MalformedFrame,
    CrcMismatch,
    ParseError,
    Unimplemented,
    ServerError,
    CommandError,
};

template <typename T>
struct Result {
    T     value;
    Error error = Error::Ok;

    bool ok() const { return error == Error::Ok; }

    static Result<T> success(T v) { return {std::move(v), Error::Ok}; }
    static Result<T> fail(Error e) { return {T{}, e}; }
};

template <>
struct Result<void> {
    Error error = Error::Ok;
    bool ok() const { return error == Error::Ok; }
    static Result<void> success() { return {Error::Ok}; }
    static Result<void> fail(Error e) { return {e}; }
};

} // namespace layrz::protocol

#endif // __LAYRZ_PROTOCOL_ERRORS_HPP__
