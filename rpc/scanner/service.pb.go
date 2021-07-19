// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.2
// source: rpc/scanner/service.proto

package scanner

import (
	common "github.com/aquasecurity/trivy/rpc/common"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ScanRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Target     string       `protobuf:"bytes,1,opt,name=target,proto3" json:"target,omitempty"` // image name or tar file path
	ArtifactId string       `protobuf:"bytes,2,opt,name=artifact_id,json=artifactId,proto3" json:"artifact_id,omitempty"`
	BlobIds    []string     `protobuf:"bytes,3,rep,name=blob_ids,json=blobIds,proto3" json:"blob_ids,omitempty"`
	Options    *ScanOptions `protobuf:"bytes,4,opt,name=options,proto3" json:"options,omitempty"`
}

func (x *ScanRequest) Reset() {
	*x = ScanRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_scanner_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScanRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScanRequest) ProtoMessage() {}

func (x *ScanRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_scanner_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScanRequest.ProtoReflect.Descriptor instead.
func (*ScanRequest) Descriptor() ([]byte, []int) {
	return file_rpc_scanner_service_proto_rawDescGZIP(), []int{0}
}

func (x *ScanRequest) GetTarget() string {
	if x != nil {
		return x.Target
	}
	return ""
}

func (x *ScanRequest) GetArtifactId() string {
	if x != nil {
		return x.ArtifactId
	}
	return ""
}

func (x *ScanRequest) GetBlobIds() []string {
	if x != nil {
		return x.BlobIds
	}
	return nil
}

func (x *ScanRequest) GetOptions() *ScanOptions {
	if x != nil {
		return x.Options
	}
	return nil
}

type ScanOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VulnType        []string `protobuf:"bytes,1,rep,name=vuln_type,json=vulnType,proto3" json:"vuln_type,omitempty"`
	SecurityChecks  []string `protobuf:"bytes,2,rep,name=security_checks,json=securityChecks,proto3" json:"security_checks,omitempty"`
	ListAllPackages bool     `protobuf:"varint,3,opt,name=list_all_packages,json=listAllPackages,proto3" json:"list_all_packages,omitempty"`
}

func (x *ScanOptions) Reset() {
	*x = ScanOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_scanner_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScanOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScanOptions) ProtoMessage() {}

func (x *ScanOptions) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_scanner_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScanOptions.ProtoReflect.Descriptor instead.
func (*ScanOptions) Descriptor() ([]byte, []int) {
	return file_rpc_scanner_service_proto_rawDescGZIP(), []int{1}
}

func (x *ScanOptions) GetVulnType() []string {
	if x != nil {
		return x.VulnType
	}
	return nil
}

func (x *ScanOptions) GetSecurityChecks() []string {
	if x != nil {
		return x.SecurityChecks
	}
	return nil
}

func (x *ScanOptions) GetListAllPackages() bool {
	if x != nil {
		return x.ListAllPackages
	}
	return false
}

type ScanResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Os      *common.OS `protobuf:"bytes,1,opt,name=os,proto3" json:"os,omitempty"`
	Eosl    bool       `protobuf:"varint,2,opt,name=eosl,proto3" json:"eosl,omitempty"`
	Results []*Result  `protobuf:"bytes,3,rep,name=results,proto3" json:"results,omitempty"`
}

func (x *ScanResponse) Reset() {
	*x = ScanResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_scanner_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScanResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScanResponse) ProtoMessage() {}

func (x *ScanResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_scanner_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScanResponse.ProtoReflect.Descriptor instead.
func (*ScanResponse) Descriptor() ([]byte, []int) {
	return file_rpc_scanner_service_proto_rawDescGZIP(), []int{2}
}

func (x *ScanResponse) GetOs() *common.OS {
	if x != nil {
		return x.Os
	}
	return nil
}

func (x *ScanResponse) GetEosl() bool {
	if x != nil {
		return x.Eosl
	}
	return false
}

func (x *ScanResponse) GetResults() []*Result {
	if x != nil {
		return x.Results
	}
	return nil
}

// Result is the same as github.com/aquasecurity/trivy/pkg/report.Result
type Result struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Target          string                  `protobuf:"bytes,1,opt,name=target,proto3" json:"target,omitempty"`
	Vulnerabilities []*common.Vulnerability `protobuf:"bytes,2,rep,name=vulnerabilities,proto3" json:"vulnerabilities,omitempty"`
	Type            string                  `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	Packages        []*common.Package       `protobuf:"bytes,4,rep,name=packages,proto3" json:"packages,omitempty"`
}

func (x *Result) Reset() {
	*x = Result{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_scanner_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Result) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Result) ProtoMessage() {}

func (x *Result) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_scanner_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Result.ProtoReflect.Descriptor instead.
func (*Result) Descriptor() ([]byte, []int) {
	return file_rpc_scanner_service_proto_rawDescGZIP(), []int{3}
}

func (x *Result) GetTarget() string {
	if x != nil {
		return x.Target
	}
	return ""
}

func (x *Result) GetVulnerabilities() []*common.Vulnerability {
	if x != nil {
		return x.Vulnerabilities
	}
	return nil
}

func (x *Result) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Result) GetPackages() []*common.Package {
	if x != nil {
		return x.Packages
	}
	return nil
}

var File_rpc_scanner_service_proto protoreflect.FileDescriptor

var file_rpc_scanner_service_proto_rawDesc = []byte{
	0x0a, 0x19, 0x72, 0x70, 0x63, 0x2f, 0x73, 0x63, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x2f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x74, 0x72, 0x69,
	0x76, 0x79, 0x2e, 0x73, 0x63, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x1a, 0x36, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x71, 0x75, 0x61, 0x73, 0x65,
	0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x2f, 0x74, 0x72, 0x69, 0x76, 0x79, 0x2f, 0x72, 0x70, 0x63,
	0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9a, 0x01, 0x0a, 0x0b, 0x53, 0x63, 0x61, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x1f, 0x0a,
	0x0b, 0x61, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x61, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x49, 0x64, 0x12, 0x19,
	0x0a, 0x08, 0x62, 0x6c, 0x6f, 0x62, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x07, 0x62, 0x6c, 0x6f, 0x62, 0x49, 0x64, 0x73, 0x12, 0x37, 0x0a, 0x07, 0x6f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x74, 0x72, 0x69,
	0x76, 0x79, 0x2e, 0x73, 0x63, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x63,
	0x61, 0x6e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x22, 0x7f, 0x0a, 0x0b, 0x53, 0x63, 0x61, 0x6e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x1b, 0x0a, 0x09, 0x76, 0x75, 0x6c, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x76, 0x75, 0x6c, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x27,
	0x0a, 0x0f, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x5f, 0x63, 0x68, 0x65, 0x63, 0x6b,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0e, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74,
	0x79, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x12, 0x2a, 0x0a, 0x11, 0x6c, 0x69, 0x73, 0x74, 0x5f,
	0x61, 0x6c, 0x6c, 0x5f, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x0f, 0x6c, 0x69, 0x73, 0x74, 0x41, 0x6c, 0x6c, 0x50, 0x61, 0x63, 0x6b, 0x61,
	0x67, 0x65, 0x73, 0x22, 0x78, 0x0a, 0x0c, 0x53, 0x63, 0x61, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x20, 0x0a, 0x02, 0x6f, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x74, 0x72, 0x69, 0x76, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4f,
	0x53, 0x52, 0x02, 0x6f, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x65, 0x6f, 0x73, 0x6c, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x04, 0x65, 0x6f, 0x73, 0x6c, 0x12, 0x32, 0x0a, 0x07, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x74, 0x72, 0x69,
	0x76, 0x79, 0x2e, 0x73, 0x63, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x52, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x22, 0xae, 0x01,
	0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x12, 0x45, 0x0a, 0x0f, 0x76, 0x75, 0x6c, 0x6e, 0x65, 0x72, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74,
	0x69, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x74, 0x72, 0x69, 0x76,
	0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x56, 0x75, 0x6c, 0x6e, 0x65, 0x72, 0x61,
	0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x52, 0x0f, 0x76, 0x75, 0x6c, 0x6e, 0x65, 0x72, 0x61, 0x62,
	0x69, 0x6c, 0x69, 0x74, 0x69, 0x65, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x31, 0x0a, 0x08, 0x70,
	0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e,
	0x74, 0x72, 0x69, 0x76, 0x79, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x50, 0x61, 0x63,
	0x6b, 0x61, 0x67, 0x65, 0x52, 0x08, 0x70, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x73, 0x32, 0x50,
	0x0a, 0x07, 0x53, 0x63, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x12, 0x45, 0x0a, 0x04, 0x53, 0x63, 0x61,
	0x6e, 0x12, 0x1d, 0x2e, 0x74, 0x72, 0x69, 0x76, 0x79, 0x2e, 0x73, 0x63, 0x61, 0x6e, 0x6e, 0x65,
	0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x63, 0x61, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1e, 0x2e, 0x74, 0x72, 0x69, 0x76, 0x79, 0x2e, 0x73, 0x63, 0x61, 0x6e, 0x6e, 0x65, 0x72,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x63, 0x61, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x42, 0x33, 0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61,
	0x71, 0x75, 0x61, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x2f, 0x74, 0x72, 0x69, 0x76,
	0x79, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x73, 0x63, 0x61, 0x6e, 0x6e, 0x65, 0x72, 0x3b, 0x73, 0x63,
	0x61, 0x6e, 0x6e, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_scanner_service_proto_rawDescOnce sync.Once
	file_rpc_scanner_service_proto_rawDescData = file_rpc_scanner_service_proto_rawDesc
)

func file_rpc_scanner_service_proto_rawDescGZIP() []byte {
	file_rpc_scanner_service_proto_rawDescOnce.Do(func() {
		file_rpc_scanner_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_scanner_service_proto_rawDescData)
	})
	return file_rpc_scanner_service_proto_rawDescData
}

var file_rpc_scanner_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_rpc_scanner_service_proto_goTypes = []interface{}{
	(*ScanRequest)(nil),          // 0: trivy.scanner.v1.ScanRequest
	(*ScanOptions)(nil),          // 1: trivy.scanner.v1.ScanOptions
	(*ScanResponse)(nil),         // 2: trivy.scanner.v1.ScanResponse
	(*Result)(nil),               // 3: trivy.scanner.v1.Result
	(*common.OS)(nil),            // 4: trivy.common.OS
	(*common.Vulnerability)(nil), // 5: trivy.common.Vulnerability
	(*common.Package)(nil),       // 6: trivy.common.Package
}
var file_rpc_scanner_service_proto_depIdxs = []int32{
	1, // 0: trivy.scanner.v1.ScanRequest.options:type_name -> trivy.scanner.v1.ScanOptions
	4, // 1: trivy.scanner.v1.ScanResponse.os:type_name -> trivy.common.OS
	3, // 2: trivy.scanner.v1.ScanResponse.results:type_name -> trivy.scanner.v1.Result
	5, // 3: trivy.scanner.v1.Result.vulnerabilities:type_name -> trivy.common.Vulnerability
	6, // 4: trivy.scanner.v1.Result.packages:type_name -> trivy.common.Package
	0, // 5: trivy.scanner.v1.Scanner.Scan:input_type -> trivy.scanner.v1.ScanRequest
	2, // 6: trivy.scanner.v1.Scanner.Scan:output_type -> trivy.scanner.v1.ScanResponse
	6, // [6:7] is the sub-list for method output_type
	5, // [5:6] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_rpc_scanner_service_proto_init() }
func file_rpc_scanner_service_proto_init() {
	if File_rpc_scanner_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_scanner_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScanRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_scanner_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScanOptions); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_scanner_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScanResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_scanner_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Result); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rpc_scanner_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_scanner_service_proto_goTypes,
		DependencyIndexes: file_rpc_scanner_service_proto_depIdxs,
		MessageInfos:      file_rpc_scanner_service_proto_msgTypes,
	}.Build()
	File_rpc_scanner_service_proto = out.File
	file_rpc_scanner_service_proto_rawDesc = nil
	file_rpc_scanner_service_proto_goTypes = nil
	file_rpc_scanner_service_proto_depIdxs = nil
}
