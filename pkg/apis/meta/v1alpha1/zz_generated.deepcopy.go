// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	net "net"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APIType) DeepCopyInto(out *APIType) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APIType.
func (in *APIType) DeepCopy() *APIType {
	if in == nil {
		return nil
	}
	out := new(APIType)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in APITypeList) DeepCopyInto(out *APITypeList) {
	{
		in := &in
		*out = make(APITypeList, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(APIType)
				(*in).DeepCopyInto(*out)
			}
		}
		return
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APITypeList.
func (in APITypeList) DeepCopy() APITypeList {
	if in == nil {
		return nil
	}
	out := new(APITypeList)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DMID) DeepCopyInto(out *DMID) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DMID.
func (in *DMID) DeepCopy() *DMID {
	if in == nil {
		return nil
	}
	out := new(DMID)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in IPAddresses) DeepCopyInto(out *IPAddresses) {
	{
		in := &in
		*out = make(IPAddresses, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = make(net.IP, len(*in))
				copy(*out, *in)
			}
		}
		return
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPAddresses.
func (in IPAddresses) DeepCopy() IPAddresses {
	if in == nil {
		return nil
	}
	out := new(IPAddresses)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectMeta) DeepCopyInto(out *ObjectMeta) {
	*out = *in
	if in.Created != nil {
		in, out := &in.Created, &out.Created
		*out = (*in).DeepCopy()
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectMeta.
func (in *ObjectMeta) DeepCopy() *ObjectMeta {
	if in == nil {
		return nil
	}
	out := new(ObjectMeta)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PortMapping) DeepCopyInto(out *PortMapping) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PortMapping.
func (in *PortMapping) DeepCopy() *PortMapping {
	if in == nil {
		return nil
	}
	out := new(PortMapping)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in PortMappings) DeepCopyInto(out *PortMappings) {
	{
		in := &in
		*out = make(PortMappings, len(*in))
		copy(*out, *in)
		return
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PortMappings.
func (in PortMappings) DeepCopy() PortMappings {
	if in == nil {
		return nil
	}
	out := new(PortMappings)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Size) DeepCopyInto(out *Size) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Size.
func (in *Size) DeepCopy() *Size {
	if in == nil {
		return nil
	}
	out := new(Size)
	in.DeepCopyInto(out)
	return out
}
