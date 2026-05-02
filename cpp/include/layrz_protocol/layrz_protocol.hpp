#pragma once

// Umbrella include — pulls in the full Layrz Protocol C++ library.
// For embedded targets only include the specific headers you need.

#include "layrz_protocol/errors.hpp"
#include "layrz_protocol/crc.hpp"
#include "layrz_protocol/float_repr.hpp"
#include "layrz_protocol/extras.hpp"
#include "layrz_protocol/parser.hpp"

#include "layrz_protocol/ble/manufacturer_data.hpp"
#include "layrz_protocol/ble/service_data.hpp"
#include "layrz_protocol/ble/advertisement.hpp"

#include "layrz_protocol/packets/firmware_branch.hpp"
#include "layrz_protocol/packets/position.hpp"
#include "layrz_protocol/packets/ble_data.hpp"
#include "layrz_protocol/packets/command.hpp"

#include "layrz_protocol/packets/ab.hpp"
#include "layrz_protocol/packets/ac.hpp"
#include "layrz_protocol/packets/ao.hpp"
#include "layrz_protocol/packets/ar.hpp"
#include "layrz_protocol/packets/as.hpp"
#include "layrz_protocol/packets/au.hpp"
#include "layrz_protocol/packets/im.hpp"
#include "layrz_protocol/packets/pa.hpp"
#include "layrz_protocol/packets/pb.hpp"
#include "layrz_protocol/packets/pc.hpp"
#include "layrz_protocol/packets/pd.hpp"
#include "layrz_protocol/packets/pi.hpp"
#include "layrz_protocol/packets/pm.hpp"
#include "layrz_protocol/packets/pr.hpp"
#include "layrz_protocol/packets/ps.hpp"
#include "layrz_protocol/packets/te.hpp"
#include "layrz_protocol/packets/ts.hpp"
