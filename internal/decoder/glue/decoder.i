%module glue
%{
#include <jxl/decode.h>

%}

%include "stdint.i"

%define JXL_EXPORT %enddef
%define JXL_DEPRECATED %enddef

%ignore JxlDecoderSetDCOutBuffer; // deprecated
%ignore JxlDecoderDCOutBufferSize; // deprecated

%include "../../include/jxl/types.h"
%include "../../include/jxl/decode.h"
%include "../../include/jxl/codestream_header.h"
