package gfx

import (
	"unsafe"

	"github.com/millken/go/sokol/util"
)

const (
	InvalidId                 = 0
	NumInflightFrames         = 2
	MaxColorAttachments       = 4
	MaxUniformBlockMembers    = 16
	MaxVertexAttributes       = 16
	MaxMipmaps                = 16
	MaxTexturearrayLayers     = 128
	MaxUniformBlockBindSlots  = 8
	MaxVertexBufferBindSlots  = 8
	MaxImageBindSlots         = 16
	MaxSamplerBindSlots       = 16
	MaxStorageBufferBindSlots = 8
	MaxImageSamplerPairs      = 16
)

type PixelFormat uint32

const (
	PixelFormatDefault PixelFormat = iota // value 0 reserved for default-init
	PixelFormatNone

	PixelFormatR8
	PixelFormatR8SN
	PixelFormatR8UI
	PixelFormatR8SI

	PixelFormatR16
	PixelFormatR16SN
	PixelFormatR16UI
	PixelFormatR16SI
	PixelFormatR16F
	PixelFormatRG8
	PixelFormatRG8SN
	PixelFormatRG8UI
	PixelFormatRG8SI

	PixelFormatR32UI
	PixelFormatR32SI
	PixelFormatR32F
	PixelFormatRG16
	PixelFormatRG16SN
	PixelFormatRG16UI
	PixelFormatRG16SI
	PixelFormatRG16F
	PixelFormatRGBA8
	PixelFormatSRGB8A8
	PixelFormatRGBA8SN
	PixelFormatRGBA8UI
	PixelFormatRGBA8SI
	PixelFormatBGRA8
	PixelFormatRGB10A2
	PixelFormatRG11B10F
	PixelFormatRGB9E5

	PixelFormatRG32UI
	PixelFormatRG32SI
	PixelFormatRG32F
	PixelFormatRGBA16
	PixelFormatRGBA16SN
	PixelFormatRGBA16UI
	PixelFormatRGBA16SI
	PixelFormatRGBA16F

	PixelFormatRGBA32UI
	PixelFormatRGBA32SI
	PixelFormatRGBA32F

	PixelFormatDepth
	PixelFormatDepthStencil

	PixelFormatBC1_RGBA
	PixelFormatBC2_RGBA
	PixelFormatBC3_RGBA
	PixelFormatBC3_SRGBA
	PixelFormatBC4_R
	PixelFormatBC4_RSN
	PixelFormatBC5_RG
	PixelFormatBC5_RGSN
	PixelFormatBC6H_RGBF
	PixelFormatBC6H_RGBUF
	PixelFormatBC7_RGBA
	PixelFormatBC7_SRGBA
	PixelFormatETC2_RGB8
	PixelFormatETC2_SRGB8
	PixelFormatETC2_RGB8A1
	PixelFormatETC2_RGBA8
	PixelFormatETC2_SRGB8A8
	PixelFormatEAC_R11
	PixelFormatEAC_R11SN
	PixelFormatEAC_RG11
	PixelFormatEAC_RG11SN

	PixelFormatASTC_4x4_RGBA
	PixelFormatASTC_4x4_SRGBA

	PixelFormatNum
)

type LoadAction uint32

const (
	LoadActionDefault LoadAction = iota
	LoadActionClear
	LoadActionLoad
	LoadActionDontCare
)

type StoreAction uint32

const (
	StoreactionDefault StoreAction = iota
	StoreactionStore
	StoreactionDontcare
)

type BufferType uint32

const (
	BufferTypeDefault BufferType = iota // value 0 reserved for default-init
	BufferTypeVertexBuffer
	BufferTypeIndexBuffer
	BufferTypeStorageBuffer
)

type IndexType uint32

const (
	IndexTypeDefault IndexType = iota // value 0 reserved for default-init
	IndexTypeNone
	IndexTypeUint16
	IndexTypeUint32
)

type Usage uint32

const (
	UsageDefault Usage = iota // value 0 reserved for default-init
	UsageImmutable
	UsageDynamic
	UsageStream
)

type ShaderAttrBaseType uint32

const (
	ShaderAttrBaseTypeUndefined ShaderAttrBaseType = iota
	ShaderAttrBaseTypeFloat
	ShaderAttrBaseTypeSint
	ShaderAttrBaseTypeUint
)

type ShaderStage uint32

const (
	ShaderStageNone ShaderStage = iota
	ShaderStageVertex
	ShaderStageFragment
	ShaderStageCompute
)

type UniformType uint32

const (
	UniformTypeInvalid UniformType = iota
	UniformTypeFloat
	UniformTypeFloat2
	UniformTypeFloat3
	UniformTypeFloat4
	UniformTypeInt
	UniformTypeInt2
	UniformTypeInt3
	UniformTypeInt4
	UniformTypeMat4
)

type UniformLayout uint32

const (
	UniformLayoutDefault UniformLayout = iota // value 0 reserved for default-init
	UniformLayoutNative                       // default: layout depends on currently active backend
	UniformLayoutStd140                       // std140: memory layout according to std140
)

type ImageType uint32

const (
	ImageTypeDefault ImageType = iota // value 0 reserved for default-init
	ImageType2d
	ImageTypeCube
	ImageType3d
	ImageTypeArray
)

type ImageSampleType uint32

const (
	ImageSampleTypeDefault ImageSampleType = iota // value 0 reserved for default-init
	ImageSampleTypeFloat
	ImageSampleTypeDepth
	ImageSampleTypeSint
	ImageSampleTypeUint
	ImageSampleTypeUnfilterableFloat
)

type SamplerType uint32

const (
	SamplerTypeDefault SamplerType = iota
	SamplerTypeFiltering
	SamplerTypeNonfiltering
	SamplerTypeComparison
)

type VertexStep uint32

const (
	VertexStepDefault VertexStep = iota // value 0 reserved for default-init
	VertexStepPerVertex
	VertexStepPerInstance
)

type VertexFormat uint32

const (
	VertexFormatInvalid VertexFormat = iota
	VertexFormatFloat
	VertexFormatFloat2
	VertexFormatFloat3
	VertexFormatFloat4
	VertexFormatInt
	VertexFormatInt2
	VertexFormatInt3
	VertexFormatInt4
	VertexFormatUint
	VertexFormatUint2
	VertexFormatUint3
	VertexFormatUint4
	VertexFormatByte4
	VertexFormatByte4N
	VertexFormatUbyte4
	VertexFormatUbyte4N
	VertexFormatShort2
	VertexFormatShort2N
	VertexFormatUshort2
	VertexFormatUshort2N
	VertexFormatShort4
	VertexFormatShort4N
	VertexFormatUshort4
	VertexFormatUshort4N
	VertexFormatUint10N2
	VertexFormatHalf2
	VertexFormatHalf4
)

type CompareFunc uint32

const (
	CompareFuncDefault CompareFunc = iota // value 0 reserved for default-init
	CompareFuncNever
	CompareFuncLess
	CompareFuncEqual
	CompareFuncLessEqual
	CompareFuncGreater
	CompareFuncNotEqual
	CompareFuncGreaterEqual
	CompareFuncAlways
)

type StencilOp uint32

const (
	StencilOpDefault StencilOp = iota // value 0 reserved for default-init
	StencilOpKeep
	StencilOpZero
	StencilOpReplace
	StencilOpIncrClamp
	StencilOpDecrClamp
	StencilOpInvert
	StencilOpIncrWrap
	StencilOpDecrWrap
)

type ColorMask uint32

const (
	ColorMaskDefault = ColorMask(0x0)  // value 0 reserved for default-init
	ColorMaskNone    = ColorMask(0x10) // special value for 'all channels disabled
	ColorMaskR       = ColorMask(0x1)
	ColorMaskG       = ColorMask(0x2)
	ColorMaskRg      = ColorMask(0x3)
	ColorMaskB       = ColorMask(0x4)
	ColorMaskRb      = ColorMask(0x5)
	ColorMaskGb      = ColorMask(0x6)
	ColorMaskRgb     = ColorMask(0x7)
	ColorMaskA       = ColorMask(0x8)
	ColorMaskRa      = ColorMask(0x9)
	ColorMaskGa      = ColorMask(0xA)
	ColorMaskRga     = ColorMask(0xB)
	ColorMaskBa      = ColorMask(0xC)
	ColorMaskRba     = ColorMask(0xD)
	ColorMaskGba     = ColorMask(0xE)
	ColorMaskRgba    = ColorMask(0xF)
)

type BlendFactor uint32

const (
	BlendFactorDefault BlendFactor = iota // value 0 reserved for default-init
	BlendFactorZero
	BlendFactorOne
	BlendFactorSrcColor
	BlendFactorOneMinusSrcColor
	BlendFactorSrcAlpha
	BlendFactorOneMinusSrcAlpha
	BlendFactorDstColor
	BlendFactorOneMinusDstColor
	BlendFactorDstAlpha
	BlendFactorOneMinusDstAlpha
	BlendFactorSrcAlphaSaturated
	BlendFactorBlendColor
	BlendFactorOneMinusBlendColor
	BlendFactorBlendAlpha
	BlendFactorOneMinusBlendAlpha
)

type BlendOp uint32

const (
	BlendOpDefault BlendOp = iota // value 0 reserved for default-init
	BlendOpAdd
	BlendOpSubtract
	BlendOpReverseSubtract
	BlendOpMin
	BlendOpMax
)

type PrimitiveType uint32

const (
	PrimitiveTypeDefault PrimitiveType = iota // value 0 reserved for default-init
	PrimitiveTypePoints
	PrimitiveTypeLines
	PrimitiveTypeLineStrip
	PrimitiveTypeTriangles
	PrimitiveTypeTriangleStrip
)

type CullMode uint32

const (
	CullModeDefault CullMode = iota // value 0 reserved for default-init
	CullModeNone
	CullModeFront
	CullModeBack
)

type FaceWinding uint32

const (
	FaceWindingDefault FaceWinding = iota // value 0 reserved for default-init
	FaceWindingCcw
	FaceWindingCw
)

type Desc struct {
	_ uint32

	BufferPoolSize                            int32
	ImagePoolSize                             int32
	SamplerPoolSize                           int32
	ShaderPoolSize                            int32
	PipelinePoolSize                          int32
	AttachmentsPoolSize                       int32
	UniformBufferSize                         int32
	MaxDispatchCallsPerPass                   int32 // max expected number of dispatch calls per pass (default: 1024)
	MaxCommitListeners                        int32
	DisableValidation                         bool  // disable validation layer even in debug mode, useful for tests
	D3D11ShaderDebugging                      bool  // if true, HLSL shaders are compiled with D3DCOMPILE_DEBUG | D3DCOMPILE_SKIP_OPTIMIZATION
	MtlForceManagedStorageMode                bool  // for debugging: use Metal managed storage mode for resources even with UMA
	MtlUseCommandBufferWithRetainedReferences bool  // Metal: use a managed MTLCommandBuffer which ref-counts used resources
	WgpuDisableBindgroupsCache                bool  // set to true to disable the WebGPU backend BindGroup cache
	WgpuBindgroupsCacheSize                   int32 // number of slots in the WebGPU bindgroup cache (must be 2^N)
	Allocator                                 Allocator
	Logger                                    Logger // optional log function override
	Environment                               Environment

	_ uint32
}

type (
	Environment struct {
		Defaults EnvironmentDefaults
		Metal    MetalEnvironment
		D3d11    D3D11Environment
		Wgpu     WgpuEnviornment
	}

	EnvironmentDefaults struct {
		ColorFormat PixelFormat
		DepthFormat PixelFormat
		SampleCount int32
	}
	MetalEnvironment struct {
		Device unsafe.Pointer
	}
	D3D11Environment struct {
		Device        unsafe.Pointer
		DeviceContext unsafe.Pointer
	}
	WgpuEnviornment struct {
		Device unsafe.Pointer
	}
)

type Logger struct {
	Func     func(tag util.CString, logLevel uint32, logItemId uint32, messageOrNil util.CString, lineNr uint32, filenameOrNil util.CString, userData unsafe.Pointer)
	UserData unsafe.Pointer
}

type Allocator struct {
	AllocFunc func(uint64, unsafe.Pointer)
	FreeFunc  func(uint64, unsafe.Pointer)
	UserData  unsafe.Pointer
}

type Color struct {
	R, G, B, A float32
}

type (
	Buffer      struct{ Id uint32 }
	Image       struct{ Id uint32 }
	Sampler     struct{ Id uint32 }
	Shader      struct{ Id uint32 }
	Pipeline    struct{ Id uint32 }
	Attachments struct{ Id uint32 }
)

type Pass struct {
	_ uint32

	Compute     bool
	Action      PassAction
	Attachments Attachments
	Swapchain   Swapchain
	Label       util.CString

	_ uint32
}

type (
	PassAction struct {
		Colors  [MaxColorAttachments]ColorAttachmentAction
		Depth   DepthAttachmentAction
		Stencil StencilAttachmentAction
	}
	ColorAttachmentAction struct {
		LoadAction  LoadAction  // default: SG_LOADACTION_CLEAR
		StoreAction StoreAction // default: SG_STOREACTION_STORE
		ClearValue  Color       // default: { 0.5f, 0.5f, 0.5f, 1.0f }
	}
	DepthAttachmentAction struct {
		LoadAction  LoadAction  // default: SG_LOADACTION_CLEAR
		StoreAction StoreAction // default: SG_STOREACTION_DONTCARE
		ClearValue  float32     // default: 1.0
	}
	StencilAttachmentAction struct {
		LoadAction  LoadAction  // default: SG_LOADACTION_CLEAR
		StoreAction StoreAction // default: SG_STOREACTION_DONTCARE
		ClearValue  uint8       // default: 0
	}
)

type (
	Swapchain struct {
		Width       int32
		Height      int32
		SampleCount int32
		ColorFormat PixelFormat
		DepthFormat PixelFormat
		Metal       MetalSwapchain
		D3D11       D3D11Swapchain
		Wgpu        WgpuSwapchain
		Gl          GlSwapchain
	}

	MetalSwapchain struct {
		CurrentDrawable     unsafe.Pointer // CAMetalDrawable (NOT MTLDrawable!!!)
		DepthStencilTexture unsafe.Pointer // MTLTexture
		MSAAColorTexture    unsafe.Pointer // MTLTexture
	}
	D3D11Swapchain struct {
		RenderView       unsafe.Pointer // ID3D11RenderTargetView
		ResolveView      unsafe.Pointer // ID3D11RenderTargetView
		DepthStencilView unsafe.Pointer // ID3D11DepthStencilView
	}
	WgpuSwapchain struct {
		RenderView       unsafe.Pointer // WGPUTextureView
		ResolveView      unsafe.Pointer // WGPUTextureView
		DepthStencilView unsafe.Pointer // WGPUTextureView
	}
	GlSwapchain struct {
		Framebuffer uint32 // GL framebuffer object
	}
)

type BufferDesc struct {
	_ uint32

	Size  uint64
	Type  BufferType
	Usage Usage
	Data  Range
	Label util.CString

	// optionally inject backend-specific resources

	GlBuffers   [NumInflightFrames]uint32
	MtlBuffers  [NumInflightFrames]unsafe.Pointer
	D3D11Buffer unsafe.Pointer
	WgpuBuffer  unsafe.Pointer

	_ uint32
}

type (
	ShaderDesc struct {
		_ uint32

		VertexFunc              ShaderFunction
		FragmentFunc            ShaderFunction
		ComputeFunc             ShaderFunction
		Attrs                   [MaxVertexAttributes]ShaderVertexAttr
		UniformBlocks           [MaxUniformBlockBindSlots]ShaderUniformBlock
		StorageBuffers          [MaxStorageBufferBindSlots]ShaderStorageBuffer
		Images                  [MaxImageBindSlots]ShaderImage
		Samplers                [MaxSamplerBindSlots]ShaderSampler
		ImageSamplerPairs       [MaxImageSamplerPairs]ShaderImageSamplerPair
		MtlTreadsPerThreadGroup MtlShaderThreadsPerThreadGroup
		Label                   util.CString

		_ uint32
	}
	ShaderFunction struct {
		Source      util.CString
		Bytecode    Range
		Entry       util.CString
		D3D11Target util.CString // default: "vs_4_0" or "ps_4_0"
	}
	ShaderVertexAttr struct {
		BaseType     ShaderAttrBaseType // default: UNDEFINED (disables validation)
		GlslName     util.CString       // [optional] GLSL attribute name
		HlslSemName  util.CString       // HLSL semantic name
		HlslSemIndex uint8              // HLSL semantic index
	}
	ShaderUniformBlock struct {
		Stage              ShaderStage
		Size               uint32
		HlslRegisterBN     uint8 // HLSL register(bn)
		MslBufferN         uint8 // MSL [[buffer(n)]]
		WgslGroup0BindingN uint8 // WGSL @group(0) @binding(n)
		Layout             UniformLayout
		GlslUniforms       [MaxUniformBlockMembers]ShaderUniform
	}
	ShaderUniform struct {
		Type       UniformType
		ArrayCount uint16       // 0 or 1 for scalars, >1 for arrays
		GlslName   util.CString // glsl name binding is required on GL 4.1 and WebGL2
	}
	ShaderStorageBuffer struct {
		Stage              ShaderStage
		Readonly           bool
		HlslRegisterTN     uint8 // HLSL register(tn) bind slot (for readonly access)
		HlslRegisterUN     uint8 // HLSL register(un) bind slot (for read/write access)
		MslBufferN         uint8 // MSL [[buffer(n)]] bind slot
		WgslGroup1BindingN uint8 // WGSL @group(1) @binding(n) bind slot
		GlslBindingN       uint8 // GLSL layout(binding=n)
	}
	ShaderImage struct {
		Stage              ShaderStage
		ImageType          ImageType
		SampleType         ImageSampleType
		Multisampled       bool
		HlslRegisterTN     uint8 // HLSL register(tn) bind slot
		MslTextureN        uint8 // MSL [[texture(n)]] bind slot
		WgslGroup1BindingN uint8 // WGSL @group(1) @binding(n) bind slot
	}
	ShaderSampler struct {
		Stage              ShaderStage
		SamplerType        SamplerType
		HlslRegisterSN     uint8 // HLSL register(sn) bind slot
		MslSamplerN        uint8 // MSL [[sampler(n)]] bind slot
		WgslGroup1BindingN uint8 // WGSL @group(1) @binding(n) bind slot
	}
	ShaderImageSamplerPair struct {
		Stage       ShaderStage
		ImageSlot   uint8
		SamplerSlot uint8
		GlslName    util.CString // glsl name binding required because of GL 4.1 and WebGL2
	}
	MtlShaderThreadsPerThreadGroup struct {
		X, Y, Z int32
	}
)

type (
	PipelineDesc struct {
		_ uint32

		Compute                bool
		Shader                 Shader
		Layout                 VertexLayoutState
		Depth                  DepthState
		Stencil                StencilState
		ColorCount             int32
		Colors                 [MaxColorAttachments]ColorTargetState
		PrimitiveType          PrimitiveType
		IndexType              IndexType
		CullMode               CullMode
		FaceWinding            FaceWinding
		SampleCount            int32
		BlendColor             Color
		AlphaToCoverageEnabled bool
		Label                  util.CString

		_ uint32
	}
	VertexLayoutState struct {
		Buffers [MaxVertexBufferBindSlots]VertexBufferLayoutState
		Attrs   [MaxVertexAttributes]VertexAttrState
	}
	VertexBufferLayoutState struct {
		Stride   int32
		StepFunc VertexStep
		StepRate int32
	}
	VertexAttrState struct {
		BufferIndex int32
		Offset      int32
		Format      VertexFormat
	}
	DepthState struct {
		PixelFormat    PixelFormat
		Compare        CompareFunc
		WriteEnabled   bool
		Bias           float32
		BiasSlopeScale float32
		BiasClamp      float32
	}
	StencilState struct {
		Enabled   bool
		Front     StencilFaceState
		Back      StencilFaceState
		ReadMask  uint8
		WriteMask uint8
		Ref       uint8
	}
	StencilFaceState struct {
		Compare     CompareFunc
		FailOp      StencilOp
		DepthFailOp StencilOp
		PassOp      StencilOp
	}
	ColorTargetState struct {
		PixelFormat PixelFormat
		WriteMask   ColorMask
		Blend       BlendState
	}
	BlendState struct {
		Enabled        bool
		SrcFactorRGB   BlendFactor
		DstFactorRgb   BlendFactor
		OpRgb          BlendOp
		SrcFactorAlpha BlendFactor
		DstFactorAlpha BlendFactor
		OpAlpha        BlendOp
	}
)

type Bindings struct {
	_ uint32

	VertexBuffers       [MaxVertexBufferBindSlots]Buffer
	VertexBufferOffsets [MaxVertexBufferBindSlots]int32
	IndexBuffer         Buffer
	IndexBufferOffset   int32
	Images              [MaxImageBindSlots]Image
	Samplers            [MaxSamplerBindSlots]Sampler
	StorageBuffers      [MaxStorageBufferBindSlots]Buffer

	_ uint32
}

type Range struct {
	Ptr  unsafe.Pointer
	Size uint64
}

func MakeRange[T any](data []T) Range {
	// @todo(Judah): If data moves at all, this will be invalidated.
	// will this even be a problem?
	return Range{
		Ptr:  unsafe.Pointer(unsafe.SliceData(data)),
		Size: uint64(unsafe.Sizeof(data[0])) * uint64(len(data)),
	}
}
