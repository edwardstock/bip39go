/*!
 * bip39.
 * bip39_config.h
 *
 * \date 2019
 * \author Eduard Maximovich (edward.vstock@gmail.com)
 * \link   https://github.com/edwardstock
 */

#ifndef BIP39_BIP39_CONFIG_H_IN
#define BIP39_BIP39_CONFIG_H_IN

#define GOEXPORT     __attribute__((visibility("default")))

#define HAVE_SYS_TYPES_H 1
#ifdef HAVE_SYS_TYPES_H
#include <sys/types.h>
#endif

#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include <stdlib.h>

typedef struct minter_data64 { uint8_t data[64]; } minter_data64;
typedef struct minter_data33 { uint8_t data[33]; } minter_data33;
typedef struct minter_data32 { uint8_t data[32]; } minter_data32;
typedef struct minter_bip32_key { uint8_t data[112]; } minter_bip32_key;

#endif //BIP39_BIP39_CONFIG_H_IN
