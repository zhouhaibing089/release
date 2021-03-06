/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// nolint: dupl
package anago_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/release/pkg/anago"
	"k8s.io/release/pkg/anago/anagofakes"
	"k8s.io/release/pkg/git"
)

func TestCheckPrerequisitesRelease(t *testing.T) {
	opts := anago.DefaultReleaseOptions()
	sut := anago.NewDefaultRelease(opts)
	mock := &anagofakes.FakeReleaseImpl{}
	sut.SetClient(mock)
	require.Nil(t, sut.CheckPrerequisites())
}

func TestGenerateReleaseVersionRelease(t *testing.T) {
	for _, tc := range []struct {
		parentBranch string
		prepare      func(*anagofakes.FakeReleaseImpl)
		shouldError  bool
	}{
		{ // success
			parentBranch: git.DefaultBranch,
			prepare:      func(*anagofakes.FakeReleaseImpl) {},
			shouldError:  false,
		},
		{ // PrepareWorkspaceRelease fails
			parentBranch: git.DefaultBranch,
			prepare: func(mock *anagofakes.FakeReleaseImpl) {
				mock.GenerateReleaseVersionReturns(nil, err)
			},
			shouldError: true,
		},
	} {
		opts := anago.DefaultReleaseOptions()
		sut := anago.NewDefaultRelease(opts)
		mock := &anagofakes.FakeReleaseImpl{}
		tc.prepare(mock)
		sut.SetClient(mock)
		_, err := sut.GenerateReleaseVersion(tc.parentBranch)
		if tc.shouldError {
			require.NotNil(t, err)
		} else {
			require.Nil(t, err)
		}
	}
}

func TestPushArtifacts(t *testing.T) {
	for _, tc := range []struct {
		prepare     func(*anagofakes.FakeReleaseImpl)
		shouldError bool
	}{
		{ // success
			prepare:     func(*anagofakes.FakeReleaseImpl) {},
			shouldError: false,
		},
		{ // CheckReleaseBucket fails
			prepare: func(mock *anagofakes.FakeReleaseImpl) {
				mock.CheckReleaseBucketReturns(err)
			},
			shouldError: true,
		},
		{ // CopyStagedFromGCSReturns fails
			prepare: func(mock *anagofakes.FakeReleaseImpl) {
				mock.CopyStagedFromGCSReturns(err)
			},
			shouldError: true,
		},
		{ // ValidateImages fails
			prepare: func(mock *anagofakes.FakeReleaseImpl) {
				mock.ValidateImagesReturns(err)
			},
			shouldError: true,
		},
		{ // PusblishVersion fails
			prepare: func(mock *anagofakes.FakeReleaseImpl) {
				mock.PublishVersionReturns(err)
			},
			shouldError: true,
		},
	} {
		opts := anago.DefaultReleaseOptions()
		sut := anago.NewDefaultRelease(opts)
		mock := &anagofakes.FakeReleaseImpl{}
		tc.prepare(mock)
		sut.SetClient(mock)
		err := sut.PushArtifacts([]string{"v1.20.0"})
		if tc.shouldError {
			require.NotNil(t, err)
		} else {
			require.Nil(t, err)
		}
	}
}

func TestPrepareWorkspaceRelease(t *testing.T) {
	for _, tc := range []struct {
		prepare     func(*anagofakes.FakeReleaseImpl)
		shouldError bool
	}{
		{ // success
			prepare:     func(*anagofakes.FakeReleaseImpl) {},
			shouldError: false,
		},
		{ // PrepareWorkspaceRelease fails
			prepare: func(mock *anagofakes.FakeReleaseImpl) {
				mock.PrepareWorkspaceReleaseReturns(err)
			},
			shouldError: true,
		},
	} {
		opts := anago.DefaultReleaseOptions()
		sut := anago.NewDefaultRelease(opts)
		mock := &anagofakes.FakeReleaseImpl{}
		tc.prepare(mock)
		sut.SetClient(mock)
		err := sut.PrepareWorkspace()
		if tc.shouldError {
			require.NotNil(t, err)
		} else {
			require.Nil(t, err)
		}
	}
}
